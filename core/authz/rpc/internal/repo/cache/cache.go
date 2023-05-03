package cache

import (
	"core/authz/proto/model"
	"core/authz/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Repository struct {
	MenuCache          *MenuCache
	RoleCache          *RoleCache
	ResourceCache      *ResourceCache
	ResourceGroupCache *ResourceGroupCache
}

func NewRepository(cfg config.Config) *Repository {
	mysqlConn := sqlx.NewMysql(cfg.MySQL.FormatDSN())
	roleCache := NewRoleCache(cfg.Cache,
		model.NewRoleModel(mysqlConn, cfg.Cache),
	)
	return &Repository{
		RoleCache: roleCache,
		MenuCache: NewMenuCache(cfg.Cache,
			model.NewMenuModel(mysqlConn, cfg.Cache),
		),
		ResourceGroupCache: NewResourceGroupCache(cfg.Cache,
			model.NewResourceGroupModel(mysqlConn),
		),
		ResourceCache: NewResourceCache(cfg.Cache, roleCache,
			model.NewResourceModel(mysqlConn, cfg.Cache),
		),
	}
}
