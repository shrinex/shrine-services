package logic

import (
	"context"

	"unit/shop/rpc/internal/svc"
	"unit/shop/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddShopConfirmLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddShopConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddShopConfirmLogic {
	return &AddShopConfirmLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddShopConfirm 创建店铺确认
func (l *AddShopConfirmLogic) AddShopConfirm(in *pb.AddShopInput) (*pb.AddShopOutput, error) {
	return &pb.AddShopOutput{}, nil
}
