package db

import (
	"core/authz/proto/model"
	"core/authz/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Repository struct {
	MenuDao          model.MenuModel
	RoleDao          model.RoleModel
	ResourceDao      model.ResourceModel
	UserRoleDao      model.UserRoleRelModel
	RoleMenuDao      model.RoleMenuRelModel
	RoleResourceDao  model.RoleResourceRelModel
	ResourceGroupDao model.ResourceGroupModel
}

func NewRepository(cfg config.Config) *Repository {
	mysqlConn := sqlx.NewMysql(cfg.MySQL.FormatDSN())
	return &Repository{
		MenuDao:          model.NewMenuModel(mysqlConn, cfg.Cache),
		RoleDao:          model.NewRoleModel(mysqlConn, cfg.Cache),
		ResourceDao:      model.NewResourceModel(mysqlConn, cfg.Cache),
		UserRoleDao:      model.NewUserRoleRelModel(mysqlConn, cfg.Cache),
		RoleMenuDao:      model.NewRoleMenuRelModel(mysqlConn, cfg.Cache),
		RoleResourceDao:  model.NewRoleResourceRelModel(mysqlConn, cfg.Cache),
		ResourceGroupDao: model.NewResourceGroupModel(mysqlConn, cfg.Cache),
	}
}
