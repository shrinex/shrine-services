package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RoleMenuRelModel = (*customRoleMenuRelModel)(nil)

type (
	// RoleMenuRelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoleMenuRelModel.
	RoleMenuRelModel interface {
		roleMenuRelModel
		DeleteByRoleId(context.Context, int64) error
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
