package model

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
	"shrine/std/utils/slices"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		UserExistsBySysTypeAndNickname(context.Context, int64, string) (bool, error)
		CountUsers(context.Context, int64, int64, string) (int64, error)
		PageUsers(context.Context, int64, int64, int64, int64, string) ([]*User, error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c),
	}
}

func (m *customUserModel) UserExistsBySysTypeAndNickname(ctx context.Context, sysType int64, nickname string) (exists bool, err error) {
	query, args := squirrel.Select("COUNT(1)").
		From("user").
		Where("sys_type = ? AND nickname = ?", sysType, nickname).
		MustSql()

	err = m.QueryRowNoCacheCtx(ctx, &exists, query, args...)
	return
}

func (m *customUserModel) CountUsers(ctx context.Context, sysType int64, shopId int64, nickname string) (int64, error) {
	builder := squirrel.Select("COUNT(1)").
		From("user").
		Where("sys_type = ? AND shop_id = ?", sysType, shopId)

	if stringx.NotEmpty(nickname) {
		builder = builder.Where(squirrel.Like{"nickname": "%" + nickname})
	}

	var resp int64
	query, args := builder.MustSql()
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

func (m *customUserModel) PageUsers(ctx context.Context, pageNo int64, pageSize int64,
	sysType int64, shopId int64, nickname string) ([]*User, error) {
	builder := squirrel.Select(userRows).
		From("user").
		Where("sys_type = ? AND shop_id = ?", sysType, shopId)

	if stringx.NotEmpty(nickname) {
		builder = builder.Where(squirrel.Like{"nickname": "%" + nickname})
	}

	query, args := builder.
		Limit(uint64(pageSize)).
		Offset(uint64((pageNo - 1) * pageSize)).
		MustSql()

	resp := slices.Empty[*User]()
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}
