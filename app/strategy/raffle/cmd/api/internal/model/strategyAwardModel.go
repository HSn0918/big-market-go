package model

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlc"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ StrategyAwardModel = (*customStrategyAwardModel)(nil)

type (
	// StrategyAwardModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStrategyAwardModel.
	StrategyAwardModel interface {
		strategyAwardModel
		QueryStrategyAwardList(ctx context.Context, StrategyId int64) (StrategyAwardList []*StrategyAward, err error)
	}

	customStrategyAwardModel struct {
		*defaultStrategyAwardModel
	}
)

// NewStrategyAwardModel returns a model for the database table.
func NewStrategyAwardModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) StrategyAwardModel {
	return &customStrategyAwardModel{
		defaultStrategyAwardModel: newStrategyAwardModel(conn, c, opts...),
	}
}
func (m *customStrategyAwardModel) QueryStrategyAwardList(ctx context.Context, StrategyId int64) (StrategyAwardList []*StrategyAward, err error) {
	//query := `SELECT strategy_id, award_id, award_title, award_subtitle, award_count, award_count_surplus, award_rate, sort FROM ` + m.table + ` WHERE strategy_id = ? ORDER BY sort ASC`
	query := `SELECT * FROM ` + m.table + ` WHERE strategy_id = ? ORDER BY sort ASC`
	err = m.QueryRowsNoCacheCtx(ctx, &StrategyAwardList, query, StrategyId)
	switch {
	case err == nil:
		return StrategyAwardList, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
