package raffle_activity

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RaffleActivityModel = (*customRaffleActivityModel)(nil)

type (
	// RaffleActivityModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRaffleActivityModel.
	RaffleActivityModel interface {
		raffleActivityModel
	}

	customRaffleActivityModel struct {
		*defaultRaffleActivityModel
	}
)

// NewRaffleActivityModel returns a model for the database table.
func NewRaffleActivityModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RaffleActivityModel {
	return &customRaffleActivityModel{
		defaultRaffleActivityModel: newRaffleActivityModel(conn, c, opts...),
	}
}
