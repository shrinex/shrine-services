package db

import (
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"unit/product/proto/model"
	"unit/product/rpc/internal/config"
)

type Repository struct {
	RawConn     sqlx.SqlConn
	CategoryDao model.CategoryModel
}

func NewRepository(cfg config.Config) *Repository {
	rawConn := sqlx.NewMysql(cfg.MySQL.FormatDSN())
	return &Repository{
		RawConn:     rawConn,
		CategoryDao: model.NewCategoryModel(rawConn, cfg.Cache),
	}
}

func (r *Repository) RawDB() *sql.DB {
	rawDB, err := r.RawConn.RawDB()
	if err != nil {
		panic(err)
	}

	return rawDB
}
