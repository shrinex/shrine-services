package db

import (
	"core/authc/proto/model"
	"core/authc/rpc/internal/config"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Repository struct {
	RawConn    sqlx.SqlConn
	UserDao    model.UserModel
	AccountDao model.AccountModel
}

func NewRepository(cfg config.Config) *Repository {
	rawConn := sqlx.NewMysql(cfg.MySQL.FormatDSN())
	return &Repository{
		RawConn:    rawConn,
		UserDao:    model.NewUserModel(rawConn, cfg.Cache),
		AccountDao: model.NewAccountModel(rawConn, cfg.Cache),
	}
}

func (r *Repository) RawDB() *sql.DB {
	rawDB, err := r.RawConn.RawDB()
	if err != nil {
		panic(err)
	}

	return rawDB
}
