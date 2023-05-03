package model

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shrine/std/utils/slices"
)

var _ RoleResourceRelModel = (*customRoleResourceRelModel)(nil)

type (
	// RoleResourceRelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoleResourceRelModel.
	RoleResourceRelModel interface {
		roleResourceRelModel
		DeleteByRoleId(context.Context, int64) error
		CountRoleIdsByResourceId(context.Context, int64) (int64, error)
		PageRoleIdsByResourceId(context.Context, int64, int64, int64) ([]int64, error)
	}

	customRoleResourceRelModel struct {
		*defaultRoleResourceRelModel
	}
)

// NewRoleResourceRelModel returns a model for the database table.
func NewRoleResourceRelModel(conn sqlx.SqlConn) RoleResourceRelModel {
	return &customRoleResourceRelModel{
		defaultRoleResourceRelModel: newRoleResourceRelModel(conn),
	}
}

func (m *customRoleResourceRelModel) DeleteByRoleId(ctx context.Context, roleId int64) error {
	query := fmt.Sprintf("delete from %s where `role_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, roleId)
	return err
}

func (m *customRoleResourceRelModel) CountRoleIdsByResourceId(ctx context.Context, resourceId int64) (int64, error) {
	builder := squirrel.Select("COUNT(1)").
		From("role_resource_rel").
		Where("resource_id = ?", resourceId)

	var resp int64
	query, args := builder.MustSql()
	err := m.conn.QueryRowCtx(ctx, &resp, query, args...)
	return resp, err
}

func (m *customRoleResourceRelModel) PageRoleIdsByResourceId(ctx context.Context, offset int64, size int64, roleId int64) ([]int64, error) {
	query, args := squirrel.Select("role_id").
		From("resource_id").
		Where("resource_id = ?", roleId).
		Limit(uint64(size)).
		Offset(uint64(offset)).
		MustSql()

	resp := slices.Empty[int64]()
	err := m.conn.QueryRowCtx(ctx, &resp, query, args...)
	return resp, err
}
