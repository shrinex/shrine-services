package model

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shrine/std/utils/slices"
)

var _ UserRoleRelModel = (*customUserRoleRelModel)(nil)

type (
	// UserRoleRelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserRoleRelModel.
	UserRoleRelModel interface {
		userRoleRelModel
		CountUserIdsByRoleId(context.Context, int64) (int64, error)
		PageUserIdsByRoleId(context.Context, int64, int64, int64) ([]int64, error)
	}

	customUserRoleRelModel struct {
		*defaultUserRoleRelModel
	}
)

// NewUserRoleRelModel returns a model for the database table.
func NewUserRoleRelModel(conn sqlx.SqlConn) UserRoleRelModel {
	return &customUserRoleRelModel{
		defaultUserRoleRelModel: newUserRoleRelModel(conn),
	}
}

func (m *customUserRoleRelModel) CountUserIdsByRoleId(ctx context.Context, roleId int64) (int64, error) {
	builder := squirrel.Select("COUNT(1)").
		From("user_role_rel").
		Where("role_id = ?", roleId)

	var resp int64
	query, args := builder.MustSql()
	err := m.conn.QueryRowCtx(ctx, &resp, query, args...)
	return resp, err
}

func (m *customUserRoleRelModel) PageUserIdsByRoleId(ctx context.Context, offset int64, size int64, roleId int64) ([]int64, error) {
	query, args := squirrel.Select("user_id").
		From("user_role_rel").
		Where("role_id = ?", roleId).
		Limit(uint64(size)).
		Offset(uint64(offset)).
		MustSql()

	resp := slices.Empty[int64]()
	err := m.conn.QueryRowCtx(ctx, &resp, query, args...)
	return resp, err
}
