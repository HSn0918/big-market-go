package svc

import (
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/config"
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/model"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                config.Config
	BizRedis              *redis.Redis
	StrategyAwardModel    model.StrategyAwardModel
	StrategyModel         model.StrategyModel
	StrategyRuleModel     model.StrategyRuleModel
	RuleTreeModel         model.RuleTreeModel
	RuleTreeNodeLineModel model.RuleTreeNodeLineModel
	RuleTreeNodeModel     model.RuleTreeNodeModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb, err := redis.NewRedis(c.BizRedis)
	if err != nil {
		logx.Error(err)
	}
	return &ServiceContext{
		Config:                c,
		BizRedis:              rdb,
		StrategyAwardModel:    model.NewStrategyAwardModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		StrategyModel:         model.NewStrategyModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		StrategyRuleModel:     model.NewStrategyRuleModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		RuleTreeModel:         model.NewRuleTreeModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		RuleTreeNodeLineModel: model.NewRuleTreeNodeLineModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		RuleTreeNodeModel:     model.NewRuleTreeNodeModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
	}
}
