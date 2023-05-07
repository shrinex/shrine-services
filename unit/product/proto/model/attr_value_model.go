package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AttrValueModel = (*customAttrValueModel)(nil)

type (
	// AttrValueModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAttrValueModel.
	AttrValueModel interface {
		attrValueModel
	}

	customAttrValueModel struct {
		*defaultAttrValueModel
	}
)

// NewAttrValueModel returns a model for the database table.
func NewAttrValueModel(conn sqlx.SqlConn, c cache.CacheConf) AttrValueModel {
	return &customAttrValueModel{
		defaultAttrValueModel: newAttrValueModel(conn, c),
	}
}
