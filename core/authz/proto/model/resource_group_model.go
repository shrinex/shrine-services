package model

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shrine/std/utils/slices"
)

var _ ResourceGroupModel = (*customResourceGroupModel)(nil)

type (
	// ResourceGroupModel is an interface to be customized, add more methods here,
	// and implement the added methods in customResourceGroupModel.
	ResourceGroupModel interface {
		resourceGroupModel
		CountResourceGroups(context.Context, int64) (int64, error)
		PageResourceGroups(context.Context, int64, int64, int64) ([]*ResourceGroup, error)
	}

	customResourceGroupModel struct {
		*defaultResourceGroupModel
	}
)

// NewResourceGroupModel returns a model for the database table.
func NewResourceGroupModel(conn sqlx.SqlConn, c cache.CacheConf) ResourceGroupModel {
	return &customResourceGroupModel{
		defaultResourceGroupModel: newResourceGroupModel(conn, c),
	}
}

func (m *customResourceGroupModel) CountResourceGroups(ctx context.Context, sysType int64) (int64, error) {
	builder := squirrel.Select("COUNT(1)").
		From("resource_group").
		Where("sys_type = ?", sysType)

	var resp int64
	query, args := builder.MustSql()
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

func (m *customResourceGroupModel) PageResourceGroups(ctx context.Context, pageNo int64, pageSize int64, sysType int64) ([]*ResourceGroup, error) {
	query, args := squirrel.Select("*").
		From("resource_group").
		Where("sys_type = ?", sysType).
		OrderBy("create_time desc").
		Limit(uint64(pageSize)).
		Offset(uint64((pageNo - 1) * pageSize)).
		MustSql()

	resp := slices.Empty[*ResourceGroup]()
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}
