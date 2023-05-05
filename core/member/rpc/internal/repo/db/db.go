package db

import (
	"core/member/proto/model"
	"core/member/rpc/internal/config"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Repository struct {
	RawConn sqlx.SqlConn
	UserDao model.UserModel
}

func NewRepository(cfg config.Config) *Repository {
	rawConn := sqlx.NewMysql(cfg.MySQL.FormatDSN())
	return &Repository{
		RawConn: rawConn,
		UserDao: model.NewUserModel(rawConn, cfg.Cache),
	}
}

func (r *Repository) RawDB() *sql.DB {
	rawDB, err := r.RawConn.RawDB()
	if err != nil {
		panic(err)
	}

	return rawDB
}
