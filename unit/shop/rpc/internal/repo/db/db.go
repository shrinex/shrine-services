package db

import (
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"unit/shop/proto/model"
	"unit/shop/rpc/internal/config"
)

type Repository struct {
	RawDB   *sql.DB
	ShopDao model.ShopModel
}

func NewRepository(cfg config.Config) *Repository {
	conn := sqlx.NewMysql(cfg.MySQL.FormatDSN())
	rawDB, _ := conn.RawDB()
	return &Repository{
		RawDB:   rawDB,
		ShopDao: model.NewShopModel(conn, cfg.Cache),
	}
}
