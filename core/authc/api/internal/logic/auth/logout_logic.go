package auth

import (
	"context"

	"core/authc/api/internal/svc"
	"core/authc/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(_ *types.LogoutReq) (resp *types.LogoutResp, err error) {
	_, err = l.svcCtx.Subject.Logout(l.ctx)
	return &types.LogoutResp{}, err
}
