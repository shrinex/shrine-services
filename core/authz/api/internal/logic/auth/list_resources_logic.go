package auth

import (
	"context"
	"core/authz/rpc/pb"
	"github.com/jinzhu/copier"
	"shrine/std/authx"
	"shrine/std/utils/slices"

	"core/authz/api/internal/svc"
	"core/authz/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListResourcesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListResourcesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListResourcesLogic {
	return &ListResourcesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListResourcesLogic) ListResources(_ *types.ListResourcesReq) (resp *types.ListResourcesResp, err error) {
	resources := slices.Empty[*types.Resource]()
	userDetails, err := l.svcCtx.Subject.UserDetails(l.ctx)
	if err != nil {
		resp = &types.ListResourcesResp{
			Resources: resources,
		}
		return
	}

	user := userDetails.(*authx.UserDetails)
	input := &pb.ListResourcesInput{
		UserId:  user.UserId,
		SysType: user.SysType,
		IsAdmin: user.IsAdmin,
	}

	output, err := l.svcCtx.AuthzRpc.ListResources(l.ctx, input)
	if err != nil {
		resp = &types.ListResourcesResp{
			Resources: resources,
		}
		return
	}

	_ = copier.Copy(&resources, output.GetResources())
	resp = &types.ListResourcesResp{
		Resources: resources,
	}
	return
}
