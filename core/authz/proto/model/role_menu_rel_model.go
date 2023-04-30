package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RoleMenuRelModel = (*customRoleMenuRelModel)(nil)

type (
	// RoleMenuRelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoleMenuRelModel.
	RoleMenuRelModel interface {
		roleMenuRelModel
	}

	customRoleMenuRelModel struct {
		*defaultRoleMenuRelModel
	}
)

// NewRoleMenuRelModel returns a model for the database table.
func NewRoleMenuRelModel(conn sqlx.SqlConn, c cache.CacheConf) RoleMenuRelModel {
	return &customRoleMenuRelModel{
		defaultRoleMenuRelModel: newRoleMenuRelModel(conn, c),
	}
}
