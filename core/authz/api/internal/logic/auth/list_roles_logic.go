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

type ListRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRolesLogic {
	return &ListRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRolesLogic) ListRoles(_ *types.ListRolesReq) (resp *types.ListRolesResp, err error) {
	roles := slices.Empty[*types.Role]()
	userDetails, err := l.svcCtx.Subject.UserDetails(l.ctx)
	if err != nil {
		resp = &types.ListRolesResp{
			Roles: roles,
		}
		return
	}

	user := userDetails.(*authx.UserDetails)
	input := &pb.ListRolesInput{
		UserId:  user.UserId,
		SysType: user.SysType,
		IsAdmin: user.IsAdmin,
	}

	output, err := l.svcCtx.AuthzRpc.ListRoles(l.ctx, input)
	if err != nil {
		resp = &types.ListRolesResp{
			Roles: roles,
		}
		return
	}

	_ = copier.Copy(&roles, output.GetRoles())
	resp = &types.ListRolesResp{
		Roles: roles,
	}
	return
}
