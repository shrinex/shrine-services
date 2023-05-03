package logic

import (
	"context"

	"core/authc/rpc/internal/svc"
	"core/authc/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterConfirmLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterConfirmLogic {
	return &RegisterConfirmLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RegisterConfirm 用户注册确认
func (l *RegisterConfirmLogic) RegisterConfirm(in *pb.RegisterInput) (*pb.RegisterOutput, error) {
	return &pb.RegisterOutput{}, nil
}
