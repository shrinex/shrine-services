package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AttrModel = (*customAttrModel)(nil)

type (
	// AttrModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAttrModel.
	AttrModel interface {
		attrModel
	}

	customAttrModel struct {
		*defaultAttrModel
	}
)

// NewAttrModel returns a model for the database table.
func NewAttrModel(conn sqlx.SqlConn, c cache.CacheConf) AttrModel {
	return &customAttrModel{
		defaultAttrModel: newAttrModel(conn, c),
	}
}
