package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BrandModel = (*customBrandModel)(nil)

type (
	// BrandModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBrandModel.
	BrandModel interface {
		brandModel
		UpdateStatus(context.Context, int64, int64) error
	}

	customBrandModel struct {
		*defaultBrandModel
	}
)

// NewBrandModel returns a model for the database table.
func NewBrandModel(conn sqlx.SqlConn, c cache.CacheConf) BrandModel {
	return &customBrandModel{
		defaultBrandModel: newBrandModel(conn, c),
	}
}

func (m *customBrandModel) UpdateStatus(ctx context.Context, brandId, status int64) error {
	brandBrandIdPrefix := fmt.Sprintf("%s%v", cacheBrandBrandIdPrefix, brandId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `status` = ? where `brand_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, status, brandId)
	}, brandBrandIdPrefix)
	return err
}
