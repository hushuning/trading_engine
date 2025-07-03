package asset

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"github.com/yzimhao/trading_engine/v2/internal/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// assetResp is the shape returned to API consumers, hiding internal fields.

type assetModule struct {
	logger    *zap.Logger
	router    *provider.Router
	assetRepo persistence.AssetRepository
}
type asset2 struct {
	id           int32
	symbol       string
	name         string
	showDecimals int
	isBase       bool
}

func newAssetModule(logger *zap.Logger, router *provider.Router, repo persistence.AssetRepository) {
	asset := assetModule{
		logger:    logger,
		router:    router,
		assetRepo: repo,
	}
	asset.registerRouter()
}

func (a *assetModule) registerRouter() {
	assetGroup := a.router.APIv1.Group("/asset")
	assetGroup.GET("/", a.query)
	assetGroup.GET("/:symbol", a.detail)

}

func (a *assetModule) query(c *gin.Context) {
	ctx := c.Request.Context()
	var list []struct {
		ID           int32  `json:"id"`
		Symbol       string `json:"symbol"`
		Name         string `json:"name"`
		ShowDecimals int    `json:"show_decimals"`
		// MinDecimals  int `json:"min_decimals"`
		IsBase bool `json:"is_base"`
		// Sort         int64 `json:"sort"`
		// Status       types.Status `json:"status"`
		// BaseAt       time.Time  `json:"base_at"`
	}
	// var symbols []string
	if err := a.assetRepo.DB().WithContext(ctx).Model(&entities.Asset{}).Select("id, symbol, name,show_decimals,is_base").Find(&list).Error; err != nil {
		a.logger.Error("list assets", zap.Error(err))
		a.router.ResponseError(c, types.ErrSystemBusy)
		return
	}
	fmt.Println(list)
	// for i := range assets {
	// 	fmt.Println(assets[i].Symbol)
	// 	fmt.Println(assets[i].Name)
	// 	fmt.Println(assets[i].ShowDecimals)
	// 	// fmt.Println(assets[i].MinDecimals)
	// 	fmt.Println(assets[i].IsBase)

	// }
	a.router.ResponseOk(c, list)
}

func (a *assetModule) detail(c *gin.Context) {
	//TODO implement
	// ctx := c.Request.Context()
	symbol := c.Param("symbol")
	if symbol == "" {
		a.router.ResponseError(c, types.ErrInvalidParam)
		return
	}

	asset, err := a.assetRepo.Get(symbol)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.router.ResponseError(c, types.ErrSystemBusy)
			return
		}
		a.logger.Error("get asset detail", zap.String("symbol", symbol), zap.Error(err))
		a.router.ResponseError(c, types.ErrSystemBusy)
		return
	}

	a.router.ResponseOk(c, asset)
}
