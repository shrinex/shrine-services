package auth

import (
	"context"
	"core/authz/rpc/pb"
	"github.com/jinzhu/copier"

	"core/authz/api/internal/svc"
	"core/authz/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveMenuLogic {
	return &RemoveMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveMenuLogic) RemoveMenu(req *types.RemoveMenuReq) (resp *types.RemoveMenuResp, err error) {
	input := new(pb.RemoveMenuInput)
	_ = copier.Copy(input, req)
	_, err = l.svcCtx.AuthzRpc.RemoveMenu(l.ctx, input)
	if err != nil {
		return
	}

	resp = &types.RemoveMenuResp{}
	return
}
