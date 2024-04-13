package test

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TestModel = (*customTestModel)(nil)

type (
	// TestModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTestModel.
	TestModel interface {
		testModel
	}

	customTestModel struct {
		*defaultTestModel
	}
)

// NewTestModel returns a model for the database table.
func NewTestModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) TestModel {
	return &customTestModel{
		defaultTestModel: newTestModel(conn, c, opts...),
	}
}
