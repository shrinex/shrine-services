package db

import (
	"core/authc/proto/model"
	"core/authc/rpc/internal/config"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Repository struct {
	RawDB      *sql.DB
	UserDao    model.UserModel
	AccountDao model.AccountModel
}

func NewRepository(cfg config.Config) *Repository {
	conn := sqlx.NewMysql(cfg.MySQL.FormatDSN())
	rawDB, _ := conn.RawDB()
	return &Repository{
		RawDB:      rawDB,
		UserDao:    model.NewUserModel(conn, cfg.Cache),
		AccountDao: model.NewAccountModel(conn, cfg.Cache),
	}
}
