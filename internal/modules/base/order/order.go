package order

import (
	"encoding/json"
	"time"

	"github.com/duolacloud/broker-core"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/modules/middlewares"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"github.com/yzimhao/trading_engine/v2/internal/types"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

// ----------------------------------
// Fx Module Registration
// ----------------------------------

var Module = fx.Module(
	"base.order",
	fx.Invoke(newOrderModule),
)

// ----------------------------------
// Structs & Constructor
// ----------------------------------

type orderModule struct {
	router *provider.Router
	logger *zap.Logger

	orderRepo persistence.OrderRepository
	broker    broker.Broker
	auth      *middlewares.AuthMiddleware
}

func newOrderModule(
	router *provider.Router,
	logger *zap.Logger,
	broker broker.Broker,
	orderRepo persistence.OrderRepository,
	auth *middlewares.AuthMiddleware, // 注入鉴权中间件
) {
	o := orderModule{
		router:    router,
		logger:    logger,
		orderRepo: orderRepo,
		broker:    broker,
		auth:      auth,
	}
	o.registerRouter()
}

// ----------------------------------
// Router Registration
// ----------------------------------

func (o *orderModule) registerRouter() {
	orderGroup := o.router.APIv1.Group("/order")

	// 权限认证：未登录将被中间件拦截并返回 401
	orderGroup.Use(o.auth.Auth())

	// POST /api/v1/order  (create order)
	orderGroup.POST("", o.create)
}

// ----------------------------------
// Request DTO
// ----------------------------------

type CreateOrderRequest struct {
	Symbol    string                   `json:"symbol" binding:"required" example:"btcusdt"`
	Side      matching_types.OrderSide `json:"side" binding:"required" example:"buy"`
	OrderType matching_types.OrderType `json:"order_type" binding:"required" example:"limit"`
	Price     *decimal.Decimal         `json:"price,omitempty" example:"1.00"`
	Quantity  *decimal.Decimal         `json:"qty,omitempty" example:"12"`
	Amount    *decimal.Decimal         `json:"amount,omitempty"`
}

// ----------------------------------
// Swagger Annotation
// ----------------------------------

// @Summary 创建订单
// @Description 创建限价 / 市价订单。限价需要 price + qty；市价至少填写 amount 或 qty。
// @ID v1.order
// @Tags order
// @Accept json
// @Produce json
// @Param args body CreateOrderRequest true "args"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/order [post]
func (o *orderModule) create(c *gin.Context) {
	var req CreateOrderRequest
	userId := o.router.ParseUserID(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		o.router.ResponseError(c, types.ErrInvalidParam)
		return
	}

	var (
		order *entities.Order
		err   error
		event types.EventOrderNew
	)

	// ---------- 限价单 ----------
	if req.OrderType == matching_types.OrderTypeLimit {
		if req.Price == nil || req.Quantity == nil {
			o.logger.Warn("price or quantity missing for limit order", zap.Any("req", req))
			o.router.ResponseError(c, types.ErrInvalidParam)
			return
		}
		order, err = o.orderRepo.CreateLimit(c, userId, req.Symbol, req.Side, *req.Price, *req.Quantity)
		if err != nil {
			o.logger.Error("create limit order", zap.Error(err), zap.Any("req", req))
			o.router.ResponseError(c, types.ErrInternalError)
			return
		}
		event.Price = &order.Price
		event.Quantity = &order.Quantity

	} else { // ---------- 市价单 ----------
		if req.Amount == nil && req.Quantity == nil {
			o.logger.Warn("amount or quantity required for market order", zap.Any("req", req))
			o.router.ResponseError(c, types.ErrInvalidParam)
			return
		}

		if req.Amount != nil && req.Amount.Cmp(decimal.Zero) > 0 {
			order, err = o.orderRepo.CreateMarketByAmount(c, userId, req.Symbol, req.Side, *req.Amount)
			if err != nil {
				o.logger.Error("create market‑amount order", zap.Error(err), zap.Any("req", req))
				o.router.ResponseError(c, types.ErrInternalError)
				return
			}
			event.Amount = &order.Amount
			event.MaxAmount = &order.FreezeAmount
		} else {
			order, err = o.orderRepo.CreateMarketByQty(c, userId, req.Symbol, req.Side, *req.Quantity)
			if err != nil {
				o.logger.Error("create market‑qty order", zap.Error(err), zap.Any("req", req))
				o.router.ResponseError(c, types.ErrInternalError)
				return
			}
			event.Quantity = &order.Quantity
			event.MaxQty = &order.FreezeQty
		}
	}

	// ---------- 发送撮合事件 ----------
	event.Symbol = order.Symbol
	event.OrderId = order.OrderId
	event.OrderSide = order.OrderSide
	event.OrderType = order.OrderType
	event.NanoTime = order.NanoTime

	body, err := json.Marshal(event)
	if err != nil {
		o.logger.Error("marshal order event", zap.Error(err), zap.Any("event", event))
		o.router.ResponseError(c, types.ErrInternalError)
		return
	}

	if err = o.broker.Publish(c, types.TOPIC_ORDER_NEW, &broker.Message{Body: body}, broker.WithShardingKey(event.Symbol)); err != nil {
		o.logger.Error("publish order event", zap.Error(err))
		o.router.ResponseError(c, types.ErrInternalError)
		return
	}

	o.router.ResponseOk(c, gin.H{"order_id": order.OrderId, "ts": time.Now().UnixNano()})
}
