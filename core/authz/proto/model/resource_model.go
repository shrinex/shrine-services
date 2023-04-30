package model

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shrine/std/utils/slices"
)

var _ ResourceModel = (*customResourceModel)(nil)

type (
	// ResourceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customResourceModel.
	ResourceModel interface {
		resourceModel
		ListAllResourcesBySysType(context.Context, int64) ([]*Resource, error)
		ListResourcesByRoleIds(context.Context, []int64) ([]*Resource, error)
		ListResourcesByRoleId(context.Context, int64) ([]*Resource, error)
	}

	customResourceModel struct {
		*defaultResourceModel
	}
)

// NewResourceModel returns a model for the database table.
func NewResourceModel(conn sqlx.SqlConn, c cache.CacheConf) ResourceModel {
	return &customResourceModel{
		defaultResourceModel: newResourceModel(conn, c),
	}
}

func (m *customResourceModel) ListAllResourcesBySysType(ctx context.Context, sysType int64) ([]*Resource, error) {
	var resp []*Resource
	query := fmt.Sprintf("select * from %s where `sys_type` = ?", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, sysType)
	if err != nil {
		return slices.Empty[*Resource](), err
	}

	return resp, nil
}

func (m *customResourceModel) ListResourcesByRoleIds(ctx context.Context, roleIds []int64) ([]*Resource, error) {
	if len(roleIds) == 0 {
		return slices.Empty[*Resource](), nil
	}

	var resp []*Resource
	query, _ := squirrel.Select("r.*").
		From("role_resource_rel rr").
		Join("resource r ON rr.resource_id = r.resource_id").
		Where(squirrel.Eq{"rr.role_id": roleIds}).
		MustSql()
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, slices.Normalize[int64](roleIds)...)
	if err != nil {
		return slices.Empty[*Resource](), err
	}

	return resp, nil
}

func (m *customResourceModel) ListResourcesByRoleId(ctx context.Context, roleId int64) ([]*Resource, error) {
	resp := slices.Empty[*Resource]()
	query, args := squirrel.Select("r.*").
		From("role_resource_rel rr").
		Join("resource r ON rr.resource_id = r.resource_id").
		Where("rr.role_id = ?", roleId).
		MustSql()
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}
