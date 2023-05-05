package db

import (
	"core/authz/proto/model"
	"core/authz/rpc/internal/config"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Repository struct {
	RawConn          sqlx.SqlConn
	MenuDao          model.MenuModel
	RoleDao          model.RoleModel
	ResourceDao      model.ResourceModel
	UserRoleDao      model.UserRoleRelModel
	RoleMenuDao      model.RoleMenuRelModel
	RoleResourceDao  model.RoleResourceRelModel
	ResourceGroupDao model.ResourceGroupModel
}

func NewRepository(cfg config.Config) *Repository {
	rawConn := sqlx.NewMysql(cfg.MySQL.FormatDSN())
	return &Repository{
		RawConn:          rawConn,
		MenuDao:          model.NewMenuModel(rawConn, cfg.Cache),
		RoleDao:          model.NewRoleModel(rawConn, cfg.Cache),
		ResourceDao:      model.NewResourceModel(rawConn, cfg.Cache),
		UserRoleDao:      model.NewUserRoleRelModel(rawConn),
		RoleMenuDao:      model.NewRoleMenuRelModel(rawConn),
		RoleResourceDao:  model.NewRoleResourceRelModel(rawConn),
		ResourceGroupDao: model.NewResourceGroupModel(rawConn),
	}
}

func (r *Repository) RawDB() *sql.DB {
	rawDB, err := r.RawConn.RawDB()
	if err != nil {
		panic(err)
	}

	return rawDB
}
