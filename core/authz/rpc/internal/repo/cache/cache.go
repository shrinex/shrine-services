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
	roleCache := NewRoleCache(cfg.Cache, model.NewRoleModel(mysqlConn, cfg.Cache))
	resourceDao := model.NewResourceModel(mysqlConn, cfg.Cache)
	menuDao := model.NewMenuModel(mysqlConn, cfg.Cache)
	resourceGroupDao := model.NewResourceGroupModel(mysqlConn, cfg.Cache)
	return &Repository{
		RoleCache:          roleCache,
		MenuCache:          NewMenuCache(cfg.Cache, menuDao),
		ResourceGroupCache: NewResourceGroupCache(cfg.Cache, resourceGroupDao),
		ResourceCache:      NewResourceCache(cfg.Cache, roleCache, resourceDao),
	}
}
