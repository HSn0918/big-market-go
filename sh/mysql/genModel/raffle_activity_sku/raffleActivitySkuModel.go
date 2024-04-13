package raffle_activity_sku

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RaffleActivitySkuModel = (*customRaffleActivitySkuModel)(nil)

type (
	// RaffleActivitySkuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRaffleActivitySkuModel.
	RaffleActivitySkuModel interface {
		raffleActivitySkuModel
	}

	customRaffleActivitySkuModel struct {
		*defaultRaffleActivitySkuModel
	}
)

// NewRaffleActivitySkuModel returns a model for the database table.
func NewRaffleActivitySkuModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RaffleActivitySkuModel {
	return &customRaffleActivitySkuModel{
		defaultRaffleActivitySkuModel: newRaffleActivitySkuModel(conn, c, opts...),
	}
}
