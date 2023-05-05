package logic

import (
	"context"

	"core/authc/rpc/internal/svc"
	"core/authc/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAdminAccountConfirmLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAdminAccountConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAdminAccountConfirmLogic {
	return &AddAdminAccountConfirmLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddAdminAccountConfirm 添加商家端管理员用户确认
func (l *AddAdminAccountConfirmLogic) AddAdminAccountConfirm(in *pb.AddAdminAccountInput) (*pb.AddAdminAccountOutput, error) {
	logx.Info("calling add admin account confirm...")
	return &pb.AddAdminAccountOutput{}, nil
}
