package model

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shrine/std/utils/slices"
)

var _ RoleMenuRelModel = (*customRoleMenuRelModel)(nil)

type (
	// RoleMenuRelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoleMenuRelModel.
	RoleMenuRelModel interface {
		roleMenuRelModel
		DeleteByRoleId(context.Context, int64) error
		CountRoleIdsByMenuId(context.Context, int64) (int64, error)
		PageRoleIdsByMenuId(context.Context, int64, int64, int64) ([]int64, error)
	}

	customRoleMenuRelModel struct {
		*defaultRoleMenuRelModel
	}
)

// NewRoleMenuRelModel returns a model for the database table.
func NewRoleMenuRelModel(conn sqlx.SqlConn) RoleMenuRelModel {
	return &customRoleMenuRelModel{
		defaultRoleMenuRelModel: newRoleMenuRelModel(conn),
	}
}

func (m *customRoleMenuRelModel) DeleteByRoleId(ctx context.Context, roleId int64) error {
	query := fmt.Sprintf("delete from %s where `role_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, roleId)
	return err
}

func (m *customRoleMenuRelModel) CountRoleIdsByMenuId(ctx context.Context, menuId int64) (int64, error) {
	builder := squirrel.Select("COUNT(1)").
		From("role_menu_rel").
		Where("menu_id = ?", menuId)

	var resp int64
	query, args := builder.MustSql()
	err := m.conn.QueryRowCtx(ctx, &resp, query, args...)
	return resp, err
}

func (m *customRoleMenuRelModel) PageRoleIdsByMenuId(ctx context.Context, offset int64, size int64, roleId int64) ([]int64, error) {
	query, args := squirrel.Select("role_id").
		From("role_menu_rel").
		Where("menu_id = ?", roleId).
		Limit(uint64(size)).
		Offset(uint64(offset)).
		MustSql()

	resp := slices.Empty[int64]()
	err := m.conn.QueryRowCtx(ctx, &resp, query, args...)
	return resp, err
}
