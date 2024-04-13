package svc

import (
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/config"
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/model"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config             config.Config
	BizRedis           *redis.Redis
	StrategyAwardModel model.StrategyAwardModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb, err := redis.NewRedis(c.BizRedis)
	if err != nil {
		logx.Error(err)
	}
	return &ServiceContext{
		Config:             c,
		BizRedis:           rdb,
		StrategyAwardModel: model.NewStrategyAwardModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
	}
}
