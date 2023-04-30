package logic

import (
	"context"

	"core/authz/rpc/internal/svc"
	"core/authz/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveRoleLogic {
	return &RemoveRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RemoveRole 删除角色
func (l *RemoveRoleLogic) RemoveRole(in *pb.RemoveRoleInput) (*pb.RemoveRoleOutput, error) {
	err := l.svcCtx.DB.RoleDao.Delete(l.ctx, in.GetRoleId())
	if err != nil {
		return nil, err
	}

	return &pb.RemoveRoleOutput{}, nil
}
