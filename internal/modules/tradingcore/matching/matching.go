package matching

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/duolacloud/broker-core"
	"github.com/duolacloud/crud-core/cache"
	"github.com/spf13/viper"
	notification_ws "github.com/yzimhao/trading_engine/v2/internal/modules/notification/ws"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	models_types "github.com/yzimhao/trading_engine/v2/internal/types"
	"github.com/yzimhao/trading_engine/v2/pkg/matching"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	CacheKeyOrderbook = "orderbook.%s" //example: orderbook.btcusdt
)

type inContext struct {
	fx.In
	Broker      broker.Broker
	Logger      *zap.Logger
	ProductRepo persistence.ProductRepository
	Viper       *viper.Viper
	Cache       cache.Cache
	OrderRepo   persistence.OrderRepository
	Ws          *notification_ws.WsManager
}

type Matching struct {
	broker      broker.Broker
	logger      *zap.Logger
	productRepo persistence.ProductRepository
	tradePairs  sync.Map
	viper       *viper.Viper
	cache       cache.Cache
	orderRepo   persistence.OrderRepository
	ws          *notification_ws.WsManager
}

func NewMatching(in inContext) *Matching {
	return &Matching{
		broker:      in.Broker,
		logger:      in.Logger,
		productRepo: in.ProductRepo,
		viper:       in.Viper,
		cache:       in.Cache,
		orderRepo:   in.OrderRepo,
		ws:          in.Ws,
	}
}

func (s *Matching) InitEngine() {
	s.logger.Sugar().Infof("init matching engine")
	localSymbols := s.viper.GetStringSlice("matching.local_symbols")

	// load trade pair

	var (
		products []entities.Product
	)

	if err := s.productRepo.DB().Model(entities.Product{}).Find(&products).Error; err != nil {
		s.logger.Sugar().Errorf("query trade product error: %v", err)
		return
	}

	for _, product := range products {
		if len(localSymbols) > 0 {
			if !slices.Contains(localSymbols, product.Symbol) {
				continue
			}
		}

		opts := []matching.Option{
			matching.WithPriceDecimals(int32(product.PriceDecimals)),
			matching.WithQuantityDecimals(int32(product.QtyDecimals)),
			// matching.WithLogger(s.logger),
		}
		engine := matching.NewEngine(context.Background(), product.Symbol, opts...)

		engine.OnRemoveResult(func(result matching_types.RemoveResult) {
			s.logger.Sugar().Infof("symbol: %s remove result: %v", result.Symbol, result)
			s.processCancelOrderResult(result)
		})
		engine.OnTradeResult(func(result matching_types.TradeResult) {
			s.logger.Sugar().Infof("symbol: %s trade result: %v", result.Symbol, result)
			s.processTradeResult(result)
		})

		s.tradePairs.Store(product.Symbol, engine)
		s.logger.Sugar().Infof("init matching engine for symbol: %s", product.Symbol)

		go s.flushOrderbookToCache(context.Background(), product.Symbol)

		//TODO  load order from db
		s.loadUnfinishedOrders(context.Background(), product.Symbol)
	}

}

func (s *Matching) Subscribe() {
	s.broker.Subscribe(models_types.TOPIC_ORDER_NEW, s.OnNewOrder)
	s.broker.Subscribe(models_types.TOPIC_NOTIFY_ORDER_CANCEL, s.OnNotifyCancelOrder)
}

func (s *Matching) OnNewOrder(ctx context.Context, event broker.Event) error {
	s.logger.Sugar().Debugf("listen new order %s", event.Message().Body)

	var order models_types.EventOrderNew
	if err := json.Unmarshal(event.Message().Body, &order); err != nil {
		s.logger.Sugar().Errorf("matching new order unmarshal error: %v body: %s", err, string(event.Message().Body))
		return err
	}

	var item matching.QueueItem
	if order.OrderType == matching_types.OrderTypeLimit {
		if order.OrderSide == matching_types.OrderSideSell {
			item = matching.NewAskLimitItem(order.OrderId, *order.Price, *order.Quantity, order.NanoTime)
		} else {
			item = matching.NewBidLimitItem(order.OrderId, *order.Price, *order.Quantity, order.NanoTime)
		}
	} else if order.OrderType == matching_types.OrderTypeMarket {
		// 按成交金额
		if order.Amount != nil {
			if order.OrderSide == matching_types.OrderSideSell {
				item = matching.NewAskMarketAmountItem(order.OrderId, *order.Amount, *order.MaxAmount, order.NanoTime)
			} else {
				item = matching.NewBidMarketAmountItem(order.OrderId, *order.Amount, order.NanoTime)
			}
		} else {
			// 按成交量
			if order.OrderSide == matching_types.OrderSideSell {
				item = matching.NewAskMarketQtyItem(order.OrderId, *order.MaxQty, order.NanoTime)
			} else {
				item = matching.NewBidMarketQtyItem(order.OrderId, *order.Quantity, *order.MaxQty, order.NanoTime)
			}
		}
	}

	if engine := s.engine(order.Symbol); engine != nil {
		engine.AddItem(item)
		s.logger.Sugar().Debugf("add item to engine %s, askLen: %d, bidLen: %d", order.Symbol, engine.AskQueue().Len(), engine.BidQueue().Len())
	}
	return nil
}

func (s *Matching) OnNotifyCancelOrder(ctx context.Context, event broker.Event) error {
	var data models_types.EventNotifyCancelOrder
	if err := json.Unmarshal(event.Message().Body, &data); err != nil {
		s.logger.Sugar().Errorf("matching notify cancel order unmarshal error: %v body: %s", err, string(event.Message().Body))
		return err
	}

	engine := s.engine(data.Symbol)
	if engine == nil {
		s.logger.Sugar().Errorf("matching engine not found for symbol: %s", data.Symbol)
		return nil
	}
	engine.RemoveItem(data.OrderSide, data.OrderId, matching_types.RemoveTypeByUser)
	return nil
}

func (s *Matching) engine(symbol string) *matching.Engine {
	if engine, ok := s.tradePairs.Load(symbol); ok {
		return engine.(*matching.Engine)
	}
	return nil
}

func (s *Matching) processCancelOrderResult(result matching_types.RemoveResult) {
	body, err := json.Marshal(result)
	if err != nil {
		s.logger.Sugar().Errorf("matching process cancel order result marshal error: %v", err)
		return
	}
	err = s.broker.Publish(context.Background(), models_types.TOPIC_PROCESS_ORDER_CANCEL, &broker.Message{
		Body: body,
	})
	if err != nil {
		s.logger.Sugar().Errorf("matching process cancel order result publish error: %v", err)
	}
}

func (s *Matching) processTradeResult(result matching_types.TradeResult) {
	body, err := result.MarshalBinary()
	if err != nil {
		s.logger.Sugar().Errorf("matching process trade result marshal error: %v", err)
		return
	}
	err = s.broker.Publish(context.Background(), models_types.TOPIC_ORDER_SETTLE, &broker.Message{
		Body: body,
	})
	if err != nil {
		s.logger.Sugar().Errorf("matching process trade result publish error: %v", err)
	}
}

func (s *Matching) flushOrderbookToCache(ctx context.Context, symbol string) {
	ticker := time.NewTicker(100 * time.Millisecond)

	engine := s.engine(symbol)
	if engine == nil {
		s.logger.Sugar().Errorf("matching engine not found for symbol: %s", symbol)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			asks := engine.GetAskOrderBook(10)
			bids := engine.GetBidOrderBook(10)

			data := map[string]any{
				"asks": asks,
				"bids": bids,
			}
			err := s.cache.Set(ctx, fmt.Sprintf(CacheKeyOrderbook, engine.Symbol()), data, cache.WithExpiration(time.Second*5))
			if err != nil {
				s.logger.Sugar().Errorf("matching flush orderbook to cache error: %v", err)
			}

			//broadcast depth data
			if err := s.ws.Broadcast(ctx, notification_ws.MsgDepthTpl.Format(map[string]string{"symbol": engine.Symbol()}), data); err != nil {
				s.logger.Sugar().Errorf("matching ws broadcast error: %v", err)
			}
		}
	}
}

func (s *Matching) loadUnfinishedOrders(ctx context.Context, symbol string) error {
	orders, err := s.orderRepo.LoadUnfinishedOrders(ctx, symbol)
	if err != nil {
		s.logger.Sugar().Errorf("matching load unfinished orders error: %v", err)
		return err
	}

	for _, order := range orders {
		var item matching.QueueItem
		if order.OrderType == matching_types.OrderTypeLimit {

			if order.OrderSide == matching_types.OrderSideSell {
				item = matching.NewAskLimitItem(order.OrderId, order.Price, order.Quantity, order.NanoTime)
			} else {
				item = matching.NewBidLimitItem(order.OrderId, order.Price, order.Quantity, order.NanoTime)
			}
			if engine := s.engine(order.Symbol); engine != nil {

				engine.AddItem(item)
				s.logger.Sugar().Debugf("load unfinished order %s to engine %s, askLen: %d, bidLen: %d", order.OrderId, order.Symbol, engine.AskQueue().Len(), engine.BidQueue().Len())
			}
		}

	}

	return nil
}
