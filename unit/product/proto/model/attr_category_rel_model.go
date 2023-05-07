package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AttrCategoryRelModel = (*customAttrCategoryRelModel)(nil)

type (
	// AttrCategoryRelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAttrCategoryRelModel.
	AttrCategoryRelModel interface {
		attrCategoryRelModel
	}

	customAttrCategoryRelModel struct {
		*defaultAttrCategoryRelModel
	}
)

// NewAttrCategoryRelModel returns a model for the database table.
func NewAttrCategoryRelModel(conn sqlx.SqlConn) AttrCategoryRelModel {
	return &customAttrCategoryRelModel{
		defaultAttrCategoryRelModel: newAttrCategoryRelModel(conn),
	}
}
