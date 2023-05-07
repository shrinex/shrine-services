package model

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shrine/std/utils/slices"
)

var _ BrandCategoryRelModel = (*customBrandCategoryRelModel)(nil)

type (
	// BrandCategoryRelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBrandCategoryRelModel.
	BrandCategoryRelModel interface {
		brandCategoryRelModel
		TxInsertBatch(context.Context, sqlx.Session, []*BrandCategoryRel) (sql.Result, error)
		TxDeleteByBrandId(context.Context, sqlx.Session, int64) error
	}

	customBrandCategoryRelModel struct {
		*defaultBrandCategoryRelModel
	}
)

// NewBrandCategoryRelModel returns a model for the database table.
func NewBrandCategoryRelModel(conn sqlx.SqlConn) BrandCategoryRelModel {
	return &customBrandCategoryRelModel{
		defaultBrandCategoryRelModel: newBrandCategoryRelModel(conn),
	}
}

func (m *customBrandCategoryRelModel) TxInsertBatch(ctx context.Context, tx sqlx.Session, rows []*BrandCategoryRel) (sql.Result, error) {
	builder := squirrel.Insert(m.table).Columns("`rel_id`", "`brand_id`", "`category_id`")
	slices.ForEach(rows, func(_ *BrandCategoryRel, _ *bool) {
		builder = builder.Values(squirrel.Expr("?, ?, ?"))
	})
	query, _ := builder.MustSql()
	args := slices.FlatMap(rows, func(e *BrandCategoryRel) []any {
		return []any{e.RelId, e.BrandId, e.CategoryId}
	})

	return tx.ExecCtx(ctx, query, args...)
}

func (m *customBrandCategoryRelModel) TxDeleteByBrandId(ctx context.Context, tx sqlx.Session, brandId int64) error {
	query, args := squirrel.Delete(m.table).Where("`brand_id` = ?", brandId).MustSql()
	_, err := tx.ExecCtx(ctx, query, args...)
	return err
}
