package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"shrine/std/conv"

	"core/authc/rpc/internal/svc"
	"core/authc/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUser 获取用户信息
func (l *GetUserLogic) GetUser(in *pb.GetUserInput) (*pb.GetUserOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	user, err := l.svcCtx.DB.UserDao.FindOne(l.ctx, in.GetUserId())
	if err != nil {
		return nil, err
	}

	var output pb.User
	_ = copier.CopyWithOption(&output, user, copier.Option{
		Converters: []copier.TypeConverter{conv.NewTime2int64()},
	})

	return &pb.GetUserOutput{
		User: &output,
	}, nil
}
