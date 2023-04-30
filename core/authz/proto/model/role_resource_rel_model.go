package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RoleResourceRelModel = (*customRoleResourceRelModel)(nil)

type (
	// RoleResourceRelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoleResourceRelModel.
	RoleResourceRelModel interface {
		roleResourceRelModel
	}

	customRoleResourceRelModel struct {
		*defaultRoleResourceRelModel
	}
)

// NewRoleResourceRelModel returns a model for the database table.
func NewRoleResourceRelModel(conn sqlx.SqlConn, c cache.CacheConf) RoleResourceRelModel {
	return &customRoleResourceRelModel{
		defaultRoleResourceRelModel: newRoleResourceRelModel(conn, c),
	}
}
