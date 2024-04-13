package model

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/hsn0918/BigMarket/common"

	"github.com/zeromicro/go-zero/core/stores/sqlc"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ StrategyRuleModel = (*customStrategyRuleModel)(nil)

type (
	// StrategyRuleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStrategyRuleModel.
	StrategyRuleModel interface {
		strategyRuleModel
		QueryStrategyRule(ctx context.Context, StrategyId int64, ruleModel string) (StrategyRule *StrategyRule, err error)
	}

	customStrategyRuleModel struct {
		*defaultStrategyRuleModel
	}
)

func (s *StrategyRule) GetRuleWeightValues() (RuleWeightValueMap map[string][]int) {
	if s.RuleModel != "rule_weight" {
		return nil
	}
	// 初始化规则权重值映射
	RuleWeightValueMap = make(map[string][]int)
	// rule_weight 权重规则配置【4000:102,103,104,105 5000:102,103,104,105,106,107 6000:102,103,104,105,106,107,108,109】
	// 1.按空格分割
	ruleGroups := strings.Split(s.RuleValue, common.SPACE)
	// 2.按冒号分割
	for _, ruleGroup := range ruleGroups {
		ruleKV := strings.Split(ruleGroup, common.COLON)
		// 3.按顿号分割
		ruleList := strings.Split(ruleKV[1], common.SPLIT)
		ruleIntList := make([]int, 0, len(ruleList))
		for _, rule := range ruleList {
			ruleInt, _ := strconv.Atoi(rule)
			ruleIntList = append(ruleIntList, ruleInt)
		}
		RuleWeightValueMap[ruleKV[0]] = append(RuleWeightValueMap[ruleKV[0]], ruleIntList...)
	}
	return RuleWeightValueMap
}

// NewStrategyRuleModel returns a model for the database table.
func NewStrategyRuleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) StrategyRuleModel {
	return &customStrategyRuleModel{
		defaultStrategyRuleModel: newStrategyRuleModel(conn, c, opts...),
	}
}
func (m *customStrategyRuleModel) QueryStrategyRule(ctx context.Context, StrategyId int64, ruleModel string) (strategyRule *StrategyRule, err error) {
	strategyRule = &StrategyRule{}
	query := `select * from` + m.table + " where strategy_id = ? and rule_model = ? limit 1"
	err = m.CachedConn.QueryRowNoCacheCtx(ctx, strategyRule, query, StrategyId, ruleModel)
	switch {
	case err == nil:
		return strategyRule, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
