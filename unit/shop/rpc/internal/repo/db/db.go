package db

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"unit/shop/proto/model"
	"unit/shop/rpc/internal/config"
)

type Repository struct {
	ShopDao model.ShopModel
}

func NewRepository(cfg config.Config) *Repository {
	return &Repository{
		ShopDao: model.NewShopModel(sqlx.NewMysql(cfg.MySQL.FormatDSN()), cfg.Cache),
	}
}
