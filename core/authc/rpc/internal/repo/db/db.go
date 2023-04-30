package db

import (
	"core/authc/proto/model"
	"core/authc/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Repository struct {
	AccountDao model.AccountModel
}

func NewRepository(cfg config.Config) *Repository {
	return &Repository{
		AccountDao: model.NewAccountModel(sqlx.NewMysql(cfg.MySQL.FormatDSN()), cfg.Cache),
	}
}
