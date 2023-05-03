package model

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AccountModel = (*customAccountModel)(nil)

type (
	// AccountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAccountModel.
	AccountModel interface {
		accountModel
		AccountExistsBySysTypeAndUsername(context.Context, int64, string) (bool, error)
	}

	customAccountModel struct {
		*defaultAccountModel
	}
)

// NewAccountModel returns a model for the database table.
func NewAccountModel(conn sqlx.SqlConn, c cache.CacheConf) AccountModel {
	return &customAccountModel{
		defaultAccountModel: newAccountModel(conn, c),
	}
}

func (m *customAccountModel) AccountExistsBySysTypeAndUsername(ctx context.Context, sysType int64, username string) (exists bool, err error) {
	query, args := squirrel.Select("COUNT(1)").
		From("account").
		Where("sys_type = ? AND username = ?", sysType, username).
		MustSql()

	err = m.QueryRowNoCacheCtx(ctx, &exists, query, args...)
	return
}
