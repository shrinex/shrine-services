package model

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shrine/std/utils/slices"
)

var _ MenuModel = (*customMenuModel)(nil)

type (
	// MenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMenuModel.
	MenuModel interface {
		menuModel
		CountMenus(context.Context, int64) (int64, error)
		PageMenus(context.Context, int64, int64, int64) ([]*Menu, error)
		ListAllMenusBySysType(context.Context, int64) ([]*Menu, error)
		ListMenusBySysTypeAndRoleId(context.Context, int64, int64) ([]*Menu, error)
	}

	customMenuModel struct {
		*defaultMenuModel
	}
)

// NewMenuModel returns a model for the database table.
func NewMenuModel(conn sqlx.SqlConn, c cache.CacheConf) MenuModel {
	return &customMenuModel{
		defaultMenuModel: newMenuModel(conn, c),
	}
}

func (m *customMenuModel) CountMenus(ctx context.Context, sysType int64) (int64, error) {
	builder := squirrel.Select("COUNT(1)").
		From("menu").
		Where("sys_type = ?", sysType)

	var resp int64
	query, args := builder.MustSql()
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

func (m *customMenuModel) PageMenus(ctx context.Context, pageNo int64, pageSize int64, sysType int64) ([]*Menu, error) {
	query, args := squirrel.Select(menuRows).
		From("menu").
		Where("sys_type = ?", sysType).
		OrderBy("weight desc").
		Limit(uint64(pageSize)).
		Offset(uint64((pageNo - 1) * pageSize)).
		MustSql()

	resp := slices.Empty[*Menu]()
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

func (m *customMenuModel) ListAllMenusBySysType(ctx context.Context, sysType int64) ([]*Menu, error) {
	var resp []*Menu
	query := fmt.Sprintf("select * from %s where `sys_type` = ?", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, sysType)
	if err != nil {
		return slices.Empty[*Menu](), err
	}

	return resp, nil
}

func (m *customMenuModel) ListMenusBySysTypeAndRoleId(ctx context.Context, sysType int64, roleId int64) ([]*Menu, error) {
	resp := slices.Empty[*Menu]()
	query, args := squirrel.Select("m.*").
		From("role_menu_rel rmr").
		Join("menu m ON rmr.menu_id = m.menu_id").
		Where("m.sys_type = ? AND rmr.role_id = ?", sysType, roleId).
		MustSql()
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}
