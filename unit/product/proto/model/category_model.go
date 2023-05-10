package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
	"shrine/std/utils/slices"
)

var _ CategoryModel = (*customCategoryModel)(nil)

type (
	// CategoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCategoryModel.
	CategoryModel interface {
		categoryModel
		ListCategories(context.Context, string) ([]*Category, error)
		UpdateStatus(context.Context, int64, int64) error
	}

	customCategoryModel struct {
		*defaultCategoryModel
	}
)

// NewCategoryModel returns a model for the database table.
func NewCategoryModel(conn sqlx.SqlConn, c cache.CacheConf) CategoryModel {
	return &customCategoryModel{
		defaultCategoryModel: newCategoryModel(conn, c),
	}
}

func (m *customCategoryModel) ListCategories(ctx context.Context, name string) ([]*Category, error) {
	builder := squirrel.Select(categoryRows).From("category")
	if stringx.NotEmpty(name) {
		builder = builder.Where(squirrel.Like{"name": name + "%"})
	}

	query, args := builder.MustSql()
	resp := slices.Empty[*Category]()
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

func (m *customCategoryModel) UpdateStatus(ctx context.Context, catId, status int64) error {
	categoryCategoryIdKey := fmt.Sprintf("%s%v", cacheCategoryCategoryIdPrefix, catId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `status` = ? where `category_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, status, catId)
	}, categoryCategoryIdKey)
	return err
}
