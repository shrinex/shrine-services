package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CategoryGroupModel = (*customCategoryGroupModel)(nil)

type (
	// CategoryGroupModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCategoryGroupModel.
	CategoryGroupModel interface {
		categoryGroupModel
	}

	customCategoryGroupModel struct {
		*defaultCategoryGroupModel
	}
)

// NewCategoryGroupModel returns a model for the database table.
func NewCategoryGroupModel(conn sqlx.SqlConn, c cache.CacheConf) CategoryGroupModel {
	return &customCategoryGroupModel{
		defaultCategoryGroupModel: newCategoryGroupModel(conn, c),
	}
}
