package auth

import (
	"context"
	"core/authz/api/internal/svc"
	"core/authz/api/internal/types"
	"core/authz/rpc/pb"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"shrine/std/authx"
	"shrine/std/globals"
)

type AddMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMenuLogic {
	return &AddMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddMenuLogic) AddMenu(req *types.AddMenuReq) (resp *types.AddMenuResp, err error) {
	resp = &types.AddMenuResp{}
	userDetails, err := l.svcCtx.Subject.UserDetails(l.ctx)
	if err != nil {
		return
	}

	user := userDetails.(*authx.UserDetails)
	if user.SysType != req.SysType {
		err = globals.ErrForbidden
		return
	}

	input := new(pb.AddMenuInput)
	_ = copier.Copy(input, req)
	output, err := l.svcCtx.AuthzRpc.AddMenu(l.ctx, input)
	if err != nil {
		return
	}

	resp.MenuId = output.GetMenuId()
	return
}
