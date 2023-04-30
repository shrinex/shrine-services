package auth

import (
	"context"
	"core/authc/api/internal/realms"
	"github.com/shrinex/shield/security"

	"core/authc/api/internal/svc"
	"core/authc/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	token := realms.NewToken(req.SysType, req.Username, req.Password)
	ctx, err := l.svcCtx.Subject.Login(l.ctx, token, security.WithRenewToken(), security.WithPlatform(req.Platform))
	if err != nil {
		resp = &types.LoginResp{}
		return
	}

	session, err := l.svcCtx.Subject.Session(ctx)
	if err != nil {
		resp = &types.LoginResp{}
		return
	}

	resp = &types.LoginResp{AccessToken: session.Token()}
	return
}
