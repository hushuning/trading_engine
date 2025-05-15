package database_test

import (
	"context"
	"testing"

	"github.com/duolacloud/crud-core/datasource"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/subosito/gotenv"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"github.com/yzimhao/trading_engine/v2/migrations"
	"go.uber.org/zap"
	_gorm "gorm.io/gorm"
)

type assetsRepoTest struct {
	suite.Suite
	ctx    context.Context
	repo   persistence.AssetRepository
	v      *viper.Viper
	gorm   *_gorm.DB
	logger *zap.Logger
}

func (suite *assetsRepoTest) SetupTest() {
	_ = gotenv.Load("../../../.env")

	suite.ctx = context.Background()

	suite.v = provider.NewViper()
	suite.gorm = provider.NewGorm(suite.v)
	suite.logger = zap.NewNop()
	redis := provider.NewRedis(suite.v, suite.logger)
	cache, _ := provider.NewCache(suite.v, redis)
	logger := zap.NewNop()
	suite.repo = database.NewAssetRepo(datasource.NewDataSource(suite.gorm), cache, logger)
}

func TestAssetsRepo(t *testing.T) {
	suite.Run(t, new(assetsRepoTest))
}

func (suite *assetsRepoTest) TearDownTest() {
	// migrations.MigrateDown(suite.gorm, suite.v, suite.logger)
}

func (suite *assetsRepoTest) TestDespoit() {
	migrations.MigrateUp(suite.gorm, suite.v, suite.logger)
	defer migrations.MigrateDown(suite.gorm, suite.v, suite.logger)

	err := suite.repo.Despoit(suite.ctx, uuid.New().String(), "user1", "BTC", types.Numeric("1"))
	suite.NoError(err)

	asset, err := suite.repo.QueryOne(suite.ctx, map[string]any{
		"symbol": map[string]any{
			"eq": "BTC",
		},
		"user_id": map[string]any{
			"eq": "user1",
		},
	})
	suite.NoError(err)
	suite.Equal("user1", asset.UserId)
	suite.Equal("BTC", asset.Symbol)
	suite.Equal(0, asset.TotalBalance.Cmp(types.Numeric("1")))
	suite.Equal(0, asset.AvailBalance.Cmp(types.Numeric("1")))
	suite.Equal(0, asset.FreezeBalance.Cmp(types.Numeric("0")))

	systemAsset, err := suite.repo.QueryOne(suite.ctx, map[string]any{
		"symbol": map[string]any{
			"eq": "BTC",
		},
		"user_id": map[string]any{
			"eq": entities.SYSTEM_USER_ROOT,
		},
	})
	suite.NoError(err)
	suite.Equal(entities.SYSTEM_USER_ROOT, systemAsset.UserId)
	suite.Equal("BTC", systemAsset.Symbol)
	suite.Equal(0, systemAsset.TotalBalance.Cmp(types.Numeric("-1")))
	suite.Equal(0, systemAsset.AvailBalance.Cmp(types.Numeric("-1")))
	suite.Equal(0, systemAsset.FreezeBalance.Cmp(types.Numeric("0")))
	//TODO test aseets_log
}

func (suite *assetsRepoTest) TestWithdraw() {

	testCases := []struct {
		name  string
		setup func()
	}{
		{
			name: "提现用户不存在",
			setup: func() {
				migrations.MigrateUp(suite.gorm, suite.v, suite.logger)
				defer migrations.MigrateDown(suite.gorm, suite.v, suite.logger)

				err := suite.repo.Withdraw(suite.ctx, uuid.New().String(), "user1", "BTC", types.Numeric("1000"))
				suite.Equal(err.Error(), "insufficient balance")
			},
		},
		{
			name: "提现用户余额不足",
			setup: func() {
				migrations.MigrateUp(suite.gorm, suite.v, suite.logger)
				defer migrations.MigrateDown(suite.gorm, suite.v, suite.logger)

				err := suite.repo.Despoit(suite.ctx, uuid.New().String(), "user1", "BTC", types.Numeric("1"))
				suite.NoError(err)

				err = suite.repo.Withdraw(suite.ctx, uuid.New().String(), "user1", "BTC", types.Numeric("1000"))
				suite.Equal(err.Error(), "insufficient balance")
			},
		},
		{
			name: "提现 余额充足",
			setup: func() {
				migrations.MigrateUp(suite.gorm, suite.v, suite.logger)
				defer migrations.MigrateDown(suite.gorm, suite.v, suite.logger)

				err := suite.repo.Despoit(suite.ctx, uuid.New().String(), "user1", "BTC", types.Numeric("2000"))
				suite.NoError(err)

				err = suite.repo.Withdraw(suite.ctx, uuid.New().String(), "user1", "BTC", types.Numeric("1000"))
				suite.NoError(err)

				asset, err := suite.repo.QueryOne(suite.ctx, map[string]any{
					"symbol": map[string]any{
						"eq": "BTC",
					},
					"user_id": map[string]any{
						"eq": "user1",
					},
				})
				// fmt.Printf("asset: %+v\n", asset)
				suite.NoError(err)
				suite.Equal("user1", asset.UserId)
				suite.Equal("BTC", asset.Symbol)
				suite.Equal(0, asset.TotalBalance.Cmp(types.Numeric("1000")))
				suite.Equal(0, asset.AvailBalance.Cmp(types.Numeric("1000")))
				suite.Equal(0, asset.FreezeBalance.Cmp(types.Numeric("0")))
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.setup()
		})
	}
}

func (suite *assetsRepoTest) TestFreeze() {
	migrations.MigrateUp(suite.gorm, suite.v, suite.logger)
	defer migrations.MigrateDown(suite.gorm, suite.v, suite.logger)

	_, err := suite.repo.Freeze(suite.ctx, suite.gorm, uuid.New().String(), "user1", "BTC", types.Numeric("1000"))
	suite.Equal(err.Error(), "insufficient balance")

	err = suite.repo.Despoit(suite.ctx, uuid.New().String(), "user1", "BTC", types.Numeric("1000"))
	suite.NoError(err)

	_, err = suite.repo.Freeze(suite.ctx, suite.gorm, uuid.New().String(), "user1", "BTC", types.Numeric("1"))
	suite.NoError(err)

	asset, err := suite.repo.QueryOne(suite.ctx, map[string]any{
		"symbol": map[string]any{
			"eq": "BTC",
		},
		"user_id": map[string]any{
			"eq": "user1",
		},
	})
	suite.NoError(err)
	suite.Equal(0, asset.FreezeBalance.Cmp(types.Numeric("1")))
	suite.Equal(0, asset.AvailBalance.Cmp(types.Numeric("999")))

	// 冻结全部
	_, err = suite.repo.Freeze(suite.ctx, suite.gorm, uuid.New().String(), "user1", "BTC", types.Numeric("0"))
	suite.NoError(err)

	asset, err = suite.repo.QueryOne(suite.ctx, map[string]any{
		"symbol": map[string]any{
			"eq": "BTC",
		},
		"user_id": map[string]any{
			"eq": "user1",
		},
	})
	suite.NoError(err)
	suite.Equal(0, asset.FreezeBalance.Cmp(types.Numeric("1000")))
	suite.Equal(0, asset.AvailBalance.Cmp(types.Numeric("0")))
}

func (suite *assetsRepoTest) TestTransfer() {
	migrations.MigrateUp(suite.gorm, suite.v, suite.logger)
	defer migrations.MigrateDown(suite.gorm, suite.v, suite.logger)

	err := suite.repo.Despoit(suite.ctx, uuid.New().String(), "user1", "BTC", types.Numeric("1000"))
	suite.NoError(err)

	transId := uuid.New().String()
	_, err = suite.repo.Freeze(suite.ctx, suite.gorm, transId, "user1", "BTC", types.Numeric("900"))
	suite.NoError(err)

	err = suite.repo.UnFreeze(suite.ctx, suite.gorm, transId, "user1", "BTC", types.Numeric("1"))
	suite.NoError(err)

	asset, err := suite.repo.QueryOne(suite.ctx, map[string]any{
		"symbol": map[string]any{
			"eq": "BTC",
		},
		"user_id": map[string]any{
			"eq": "user1",
		},
	})
	suite.NoError(err)
	suite.Equal(0, asset.FreezeBalance.Cmp(types.Numeric("899")))
	suite.Equal(0, asset.AvailBalance.Cmp(types.Numeric("101")))

	//解冻全部
	err = suite.repo.UnFreeze(suite.ctx, suite.gorm, transId, "user1", "BTC", types.Numeric("0"))
	suite.NoError(err)

	asset, err = suite.repo.QueryOne(suite.ctx, map[string]any{
		"symbol": map[string]any{
			"eq": "BTC",
		},
		"user_id": map[string]any{
			"eq": "user1",
		},
	})
	suite.NoError(err)
	suite.Equal(0, asset.FreezeBalance.Cmp(types.Numeric("0")))
	suite.Equal(0, asset.AvailBalance.Cmp(types.Numeric("1000")))
}
