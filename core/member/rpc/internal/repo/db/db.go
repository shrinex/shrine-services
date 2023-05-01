package db

import (
	"core/member/proto/model"
	"core/member/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Repository struct {
	UserDao model.UserModel
}

func NewRepository(cfg config.Config) *Repository {
	return &Repository{
		UserDao: model.NewUserModel(sqlx.NewMysql(cfg.MySQL.FormatDSN()), cfg.Cache),
	}
}
