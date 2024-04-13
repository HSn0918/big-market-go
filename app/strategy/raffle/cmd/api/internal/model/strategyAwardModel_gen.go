// Code generated by goctl. DO NOT EDIT.

package model

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
	strategyAwardFieldNames          = builder.RawFieldNames(&StrategyAward{})
	strategyAwardRows                = strings.Join(strategyAwardFieldNames, ",")
	strategyAwardRowsExpectAutoSet   = strings.Join(stringx.Remove(strategyAwardFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	strategyAwardRowsWithPlaceHolder = strings.Join(stringx.Remove(strategyAwardFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheBigMarketStrategyAwardIdPrefix = "cache:bigMarket:strategyAward:id:"
)

type (
	strategyAwardModel interface {
		Insert(ctx context.Context, data *StrategyAward) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*StrategyAward, error)
		Update(ctx context.Context, data *StrategyAward) error
		Delete(ctx context.Context, id int64) error
	}

	defaultStrategyAwardModel struct {
		sqlc.CachedConn
		table string
	}

	StrategyAward struct {
		Id                int64          `db:"id"`                  // 自增ID
		StrategyId        int64          `db:"strategy_id"`         // 抽奖策略ID
		AwardId           int64          `db:"award_id"`            // 抽奖奖品ID - 内部流转使用
		AwardTitle        string         `db:"award_title"`         // 抽奖奖品标题
		AwardSubtitle     sql.NullString `db:"award_subtitle"`      // 抽奖奖品副标题
		AwardCount        int64          `db:"award_count"`         // 奖品库存总量
		AwardCountSurplus int64          `db:"award_count_surplus"` // 奖品库存剩余
		AwardRate         float64        `db:"award_rate"`          // 奖品中奖概率
		RuleModels        sql.NullString `db:"rule_models"`         // 规则模型，rule配置的模型同步到此表，便于使用
		Sort              int64          `db:"sort"`                // 排序
		CreateTime        time.Time      `db:"create_time"`         // 创建时间
		UpdateTime        time.Time      `db:"update_time"`         // 修改时间
	}
)

func newStrategyAwardModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultStrategyAwardModel {
	return &defaultStrategyAwardModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`strategy_award`",
	}
}

func (m *defaultStrategyAwardModel) Delete(ctx context.Context, id int64) error {
	bigMarketStrategyAwardIdKey := fmt.Sprintf("%s%v", cacheBigMarketStrategyAwardIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, bigMarketStrategyAwardIdKey)
	return err
}

func (m *defaultStrategyAwardModel) FindOne(ctx context.Context, id int64) (*StrategyAward, error) {
	bigMarketStrategyAwardIdKey := fmt.Sprintf("%s%v", cacheBigMarketStrategyAwardIdPrefix, id)
	var resp StrategyAward
	err := m.QueryRowCtx(ctx, &resp, bigMarketStrategyAwardIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", strategyAwardRows, m.table)
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

func (m *defaultStrategyAwardModel) Insert(ctx context.Context, data *StrategyAward) (sql.Result, error) {
	bigMarketStrategyAwardIdKey := fmt.Sprintf("%s%v", cacheBigMarketStrategyAwardIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, strategyAwardRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.StrategyId, data.AwardId, data.AwardTitle, data.AwardSubtitle, data.AwardCount, data.AwardCountSurplus, data.AwardRate, data.RuleModels, data.Sort)
	}, bigMarketStrategyAwardIdKey)
	return ret, err
}

func (m *defaultStrategyAwardModel) Update(ctx context.Context, data *StrategyAward) error {
	bigMarketStrategyAwardIdKey := fmt.Sprintf("%s%v", cacheBigMarketStrategyAwardIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strategyAwardRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.StrategyId, data.AwardId, data.AwardTitle, data.AwardSubtitle, data.AwardCount, data.AwardCountSurplus, data.AwardRate, data.RuleModels, data.Sort, data.Id)
	}, bigMarketStrategyAwardIdKey)
	return err
}

func (m *defaultStrategyAwardModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheBigMarketStrategyAwardIdPrefix, primary)
}

func (m *defaultStrategyAwardModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", strategyAwardRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultStrategyAwardModel) tableName() string {
	return m.table
}