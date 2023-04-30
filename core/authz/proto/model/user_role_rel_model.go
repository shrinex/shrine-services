package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserRoleRelModel = (*customUserRoleRelModel)(nil)

type (
	// UserRoleRelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserRoleRelModel.
	UserRoleRelModel interface {
		userRoleRelModel
	}

	customUserRoleRelModel struct {
		*defaultUserRoleRelModel
	}
)

// NewUserRoleRelModel returns a model for the database table.
func NewUserRoleRelModel(conn sqlx.SqlConn, c cache.CacheConf) UserRoleRelModel {
	return &customUserRoleRelModel{
		defaultUserRoleRelModel: newUserRoleRelModel(conn, c),
	}
}
