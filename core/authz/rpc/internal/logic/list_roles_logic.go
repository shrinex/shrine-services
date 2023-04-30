package logic

import (
	"context"
	"core/authz/rpc/internal/repo/cache"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"shrine/std/conv"

	"core/authz/rpc/internal/svc"
	"core/authz/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRolesLogic {
	return &ListRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListRoles 查询用户拥有的角色列表
func (l *ListRolesLogic) ListRoles(in *pb.ListRolesInput) (*pb.ListRolesOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	roles, err := l.svcCtx.Cache.RoleCache.ListRoles(l.ctx, in.GetSysType(), in.GetUserId(), in.GetIsAdmin())
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		return nil, err
	}

	var outputs []*pb.Role
	_ = copier.CopyWithOption(&outputs, roles, copier.Option{
		Converters: []copier.TypeConverter{conv.NewTime2int64()},
	})

	return &pb.ListRolesOutput{
		Roles: outputs,
	}, nil
}
