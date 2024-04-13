package raffle_activity_count

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RaffleActivityCountModel = (*customRaffleActivityCountModel)(nil)

type (
	// RaffleActivityCountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRaffleActivityCountModel.
	RaffleActivityCountModel interface {
		raffleActivityCountModel
	}

	customRaffleActivityCountModel struct {
		*defaultRaffleActivityCountModel
	}
)

// NewRaffleActivityCountModel returns a model for the database table.
func NewRaffleActivityCountModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RaffleActivityCountModel {
	return &customRaffleActivityCountModel{
		defaultRaffleActivityCountModel: newRaffleActivityCountModel(conn, c, opts...),
	}
}
