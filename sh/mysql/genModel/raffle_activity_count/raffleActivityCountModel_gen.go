// Code generated by goctl. DO NOT EDIT.

package raffle_activity_count

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
	raffleActivityCountFieldNames          = builder.RawFieldNames(&RaffleActivityCount{})
	raffleActivityCountRows                = strings.Join(raffleActivityCountFieldNames, ",")
	raffleActivityCountRowsExpectAutoSet   = strings.Join(stringx.Remove(raffleActivityCountFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	raffleActivityCountRowsWithPlaceHolder = strings.Join(stringx.Remove(raffleActivityCountFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheBigMarketRaffleActivityCountIdPrefix              = "cache:bigMarket:raffleActivityCount:id:"
	cacheBigMarketRaffleActivityCountActivityCountIdPrefix = "cache:bigMarket:raffleActivityCount:activityCountId:"
)

type (
	raffleActivityCountModel interface {
		Insert(ctx context.Context, data *RaffleActivityCount) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*RaffleActivityCount, error)
		FindOneByActivityCountId(ctx context.Context, activityCountId int64) (*RaffleActivityCount, error)
		Update(ctx context.Context, data *RaffleActivityCount) error
		Delete(ctx context.Context, id int64) error
	}

	defaultRaffleActivityCountModel struct {
		sqlc.CachedConn
		table string
	}

	RaffleActivityCount struct {
		Id              int64     `db:"id"`                // 自增ID
		ActivityCountId int64     `db:"activity_count_id"` // 活动次数编号
		TotalCount      int64     `db:"total_count"`       // 总次数
		DayCount        int64     `db:"day_count"`         // 日次数
		MonthCount      int64     `db:"month_count"`       // 月次数
		CreateTime      time.Time `db:"create_time"`       // 创建时间
		UpdateTime      time.Time `db:"update_time"`       // 更新时间
	}
)

func newRaffleActivityCountModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultRaffleActivityCountModel {
	return &defaultRaffleActivityCountModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`raffle_activity_count`",
	}
}

func (m *defaultRaffleActivityCountModel) Delete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	bigMarketRaffleActivityCountActivityCountIdKey := fmt.Sprintf("%s%v", cacheBigMarketRaffleActivityCountActivityCountIdPrefix, data.ActivityCountId)
	bigMarketRaffleActivityCountIdKey := fmt.Sprintf("%s%v", cacheBigMarketRaffleActivityCountIdPrefix, id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, bigMarketRaffleActivityCountActivityCountIdKey, bigMarketRaffleActivityCountIdKey)
	return err
}

func (m *defaultRaffleActivityCountModel) FindOne(ctx context.Context, id int64) (*RaffleActivityCount, error) {
	bigMarketRaffleActivityCountIdKey := fmt.Sprintf("%s%v", cacheBigMarketRaffleActivityCountIdPrefix, id)
	var resp RaffleActivityCount
	err := m.QueryRowCtx(ctx, &resp, bigMarketRaffleActivityCountIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", raffleActivityCountRows, m.table)
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

func (m *defaultRaffleActivityCountModel) FindOneByActivityCountId(ctx context.Context, activityCountId int64) (*RaffleActivityCount, error) {
	bigMarketRaffleActivityCountActivityCountIdKey := fmt.Sprintf("%s%v", cacheBigMarketRaffleActivityCountActivityCountIdPrefix, activityCountId)
	var resp RaffleActivityCount
	err := m.QueryRowIndexCtx(ctx, &resp, bigMarketRaffleActivityCountActivityCountIdKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where `activity_count_id` = ? limit 1", raffleActivityCountRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, activityCountId); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultRaffleActivityCountModel) Insert(ctx context.Context, data *RaffleActivityCount) (sql.Result, error) {
	bigMarketRaffleActivityCountActivityCountIdKey := fmt.Sprintf("%s%v", cacheBigMarketRaffleActivityCountActivityCountIdPrefix, data.ActivityCountId)
	bigMarketRaffleActivityCountIdKey := fmt.Sprintf("%s%v", cacheBigMarketRaffleActivityCountIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, raffleActivityCountRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.ActivityCountId, data.TotalCount, data.DayCount, data.MonthCount)
	}, bigMarketRaffleActivityCountActivityCountIdKey, bigMarketRaffleActivityCountIdKey)
	return ret, err
}

func (m *defaultRaffleActivityCountModel) Update(ctx context.Context, newData *RaffleActivityCount) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	bigMarketRaffleActivityCountActivityCountIdKey := fmt.Sprintf("%s%v", cacheBigMarketRaffleActivityCountActivityCountIdPrefix, data.ActivityCountId)
	bigMarketRaffleActivityCountIdKey := fmt.Sprintf("%s%v", cacheBigMarketRaffleActivityCountIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, raffleActivityCountRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.ActivityCountId, newData.TotalCount, newData.DayCount, newData.MonthCount, newData.Id)
	}, bigMarketRaffleActivityCountActivityCountIdKey, bigMarketRaffleActivityCountIdKey)
	return err
}

func (m *defaultRaffleActivityCountModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheBigMarketRaffleActivityCountIdPrefix, primary)
}

func (m *defaultRaffleActivityCountModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", raffleActivityCountRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultRaffleActivityCountModel) tableName() string {
	return m.table
}
