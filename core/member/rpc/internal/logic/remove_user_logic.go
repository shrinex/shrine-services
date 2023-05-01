package logic

import (
	"context"

	"core/member/rpc/internal/svc"
	"core/member/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveUserLogic {
	return &RemoveUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RemoveUser 删除用户
func (l *RemoveUserLogic) RemoveUser(in *pb.RemoveUserInput) (*pb.RemoveUserOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.DB.UserDao.Delete(l.ctx, in.GetUserId())
	if err != nil {
		return nil, err
	}

	return &pb.RemoveUserOutput{}, nil
}
