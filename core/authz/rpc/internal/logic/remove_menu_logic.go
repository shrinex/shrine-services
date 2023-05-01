package logic

import (
	"context"

	"core/authz/rpc/internal/svc"
	"core/authz/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveMenuLogic {
	return &RemoveMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RemoveMenu 删除菜单
func (l *RemoveMenuLogic) RemoveMenu(in *pb.RemoveMenuInput) (*pb.RemoveMenuOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.DB.MenuDao.Delete(l.ctx, in.GetMenuId())
	if err != nil {
		return nil, err
	}

	return &pb.RemoveMenuOutput{}, nil
}
