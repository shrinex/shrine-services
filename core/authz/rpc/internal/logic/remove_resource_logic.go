package logic

import (
	"context"

	"core/authz/rpc/internal/svc"
	"core/authz/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveResourceLogic {
	return &RemoveResourceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RemoveResource 删除资源
func (l *RemoveResourceLogic) RemoveResource(in *pb.RemoveResourceInput) (*pb.RemoveResourceOutput, error) {
	err := l.svcCtx.DB.ResourceDao.Delete(l.ctx, in.GetResourceId())
	if err != nil {
		return nil, err
	}

	return &pb.RemoveResourceOutput{}, nil
}
