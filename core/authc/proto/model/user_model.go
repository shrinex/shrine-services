package model

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		UserExistsBySysTypeAndNickname(context.Context, int64, string) (bool, error)
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
