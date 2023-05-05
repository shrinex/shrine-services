package db

import (
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"unit/shop/proto/model"
	"unit/shop/rpc/internal/config"
)

type Repository struct {
	RawConn sqlx.SqlConn
	ShopDao model.ShopModel
}

func NewRepository(cfg config.Config) *Repository {
	rawConn := sqlx.NewMysql(cfg.MySQL.FormatDSN())
	return &Repository{
		RawConn: rawConn,
		ShopDao: model.NewShopModel(rawConn, cfg.Cache),
	}
}

func (r *Repository) RawDB() *sql.DB {
	rawDB, err := r.RawConn.RawDB()
	if err != nil {
		panic(err)
	}

	return rawDB
}
