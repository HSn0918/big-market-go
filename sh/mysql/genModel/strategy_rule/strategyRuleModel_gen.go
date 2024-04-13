// Code generated by goctl. DO NOT EDIT.

package strategy_rule

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	strategyRuleFieldNames          = builder.RawFieldNames(&StrategyRule{})
	strategyRuleRows                = strings.Join(strategyRuleFieldNames, ",")
	strategyRuleRowsExpectAutoSet   = strings.Join(stringx.Remove(strategyRuleFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	strategyRuleRowsWithPlaceHolder = strings.Join(stringx.Remove(strategyRuleFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheBigMarketStrategyRuleIdPrefix = "cache:bigMarket:strategyRule:id:"
)

type (
	strategyRuleModel interface {
		Insert(ctx context.Context, data *StrategyRule) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*StrategyRule, error)
		Update(ctx context.Context, data *StrategyRule) error
		Delete(ctx context.Context, id int64) error
	}

	defaultStrategyRuleModel struct {
		sqlc.CachedConn
		table string
	}

	StrategyRule struct {
		Id         int64         `db:"id"`          // 自增ID
		StrategyId int64         `db:"strategy_id"` // 抽奖策略ID
		AwardId    sql.NullInt64 `db:"award_id"`    // 抽奖奖品ID【规则类型为策略，则不需要奖品ID】
		RuleType   int64         `db:"rule_type"`   // 抽象规则类型；1-策略规则、2-奖品规则
		RuleModel  string        `db:"rule_model"`  // 抽奖规则类型【rule_random - 随机值计算、rule_lock - 抽奖几次后解锁、rule_luck_award - 幸运奖(兜底奖品)】
		RuleValue  string        `db:"rule_value"`  // 抽奖规则比值
		RuleDesc   string        `db:"rule_desc"`   // 抽奖规则描述
		CreateTime time.Time     `db:"create_time"` // 创建时间
		UpdateTime time.Time     `db:"update_time"` // 更新时间
	}
)

func newStrategyRuleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultStrategyRuleModel {
	return &defaultStrategyRuleModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`strategy_rule`",
	}
}

func (m *defaultStrategyRuleModel) Delete(ctx context.Context, id int64) error {
	bigMarketStrategyRuleIdKey := fmt.Sprintf("%s%v", cacheBigMarketStrategyRuleIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, bigMarketStrategyRuleIdKey)
	return err
}

func (m *defaultStrategyRuleModel) FindOne(ctx context.Context, id int64) (*StrategyRule, error) {
	bigMarketStrategyRuleIdKey := fmt.Sprintf("%s%v", cacheBigMarketStrategyRuleIdPrefix, id)
	var resp StrategyRule
	err := m.QueryRowCtx(ctx, &resp, bigMarketStrategyRuleIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", strategyRuleRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultStrategyRuleModel) Insert(ctx context.Context, data *StrategyRule) (sql.Result, error) {
	bigMarketStrategyRuleIdKey := fmt.Sprintf("%s%v", cacheBigMarketStrategyRuleIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, strategyRuleRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.StrategyId, data.AwardId, data.RuleType, data.RuleModel, data.RuleValue, data.RuleDesc)
	}, bigMarketStrategyRuleIdKey)
	return ret, err
}

func (m *defaultStrategyRuleModel) Update(ctx context.Context, data *StrategyRule) error {
	bigMarketStrategyRuleIdKey := fmt.Sprintf("%s%v", cacheBigMarketStrategyRuleIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strategyRuleRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.StrategyId, data.AwardId, data.RuleType, data.RuleModel, data.RuleValue, data.RuleDesc, data.Id)
	}, bigMarketStrategyRuleIdKey)
	return err
}

func (m *defaultStrategyRuleModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheBigMarketStrategyRuleIdPrefix, primary)
}

func (m *defaultStrategyRuleModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", strategyRuleRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultStrategyRuleModel) tableName() string {
	return m.table
}
