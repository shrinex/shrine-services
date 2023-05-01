package logic

import (
	"context"

	"core/member/rpc/internal/svc"
	"core/member/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserExistsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserExistsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserExistsLogic {
	return &UserExistsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserExists 判断用户是否已存在
func (l *UserExistsLogic) UserExists(in *pb.UserExistsInput) (*pb.UserExistsOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	exists, err := l.svcCtx.DB.UserDao.UserExistsBySysTypeAndNickname(l.ctx, in.GetSysType(), in.GetNickname())
	if err != nil {
		return nil, err
	}

	return &pb.UserExistsOutput{
		Exists: exists,
	}, nil
}
