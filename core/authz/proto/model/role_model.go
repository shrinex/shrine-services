package model

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shrine/std/utils/slices"
)

var _ RoleModel = (*customRoleModel)(nil)

type (
	// RoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoleModel.
	RoleModel interface {
		roleModel
		ListAllRolesBySysType(context.Context, int64) ([]*Role, error)
		ListRolesByUserIdAndSysType(context.Context, int64, int64) ([]*Role, error)
	}

	customRoleModel struct {
		*defaultRoleModel
	}
)

// NewRoleModel returns a model for the database table.
func NewRoleModel(conn sqlx.SqlConn, c cache.CacheConf) RoleModel {
	return &customRoleModel{
		defaultRoleModel: newRoleModel(conn, c),
	}
}

func (m *customRoleModel) ListAllRolesBySysType(ctx context.Context, sysType int64) ([]*Role, error) {
	var resp []*Role
	query := fmt.Sprintf("select * from %s where `sys_type` = ?", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, sysType)
	if err != nil {
		return slices.Empty[*Role](), err
	}

	return resp, nil
}

func (m *customRoleModel) ListRolesByUserIdAndSysType(ctx context.Context, userId int64, sysType int64) ([]*Role, error) {
	var resp []*Role
	query, args := squirrel.Select("r.*").
		From("user_role_rel ur").
		Join("role r ON ur.role_id = r.role_id").
		Where("ur.user_id = ?  AND r.sys_type = ?", userId, sysType).
		MustSql()
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	if err != nil {
		return slices.Empty[*Role](), err
	}

	return resp, nil
}
