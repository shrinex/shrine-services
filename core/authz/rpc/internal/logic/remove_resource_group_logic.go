package logic

import (
	"context"

	"core/authz/rpc/internal/svc"
	"core/authz/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveResourceGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveResourceGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveResourceGroupLogic {
	return &RemoveResourceGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RemoveResourceGroup 删除资源分组
func (l *RemoveResourceGroupLogic) RemoveResourceGroup(in *pb.RemoveResourceGroupInput) (*pb.RemoveResourceGroupOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.DB.ResourceGroupDao.Delete(l.ctx, in.GetGroupId())
	if err != nil {
		return nil, err
	}

	return &pb.RemoveResourceGroupOutput{}, nil
}
