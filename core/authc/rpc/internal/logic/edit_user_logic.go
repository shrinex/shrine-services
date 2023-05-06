package logic

import (
	"context"
	"core/authc/proto/model"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"core/authc/rpc/internal/svc"
	"core/authc/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditUserLogic {
	return &EditUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// EditUser 编辑用户信息
func (l *EditUserLogic) EditUser(in *pb.EditUserInput) (*pb.EditUserOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	exist, err := l.validate(in)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		UserId:   exist.UserId,
		ShopId:   exist.ShopId,
		SysType:  exist.SysType,
		Nickname: in.GetNickname(),
		Avatar:   in.GetAvatar(),
		Intro:    in.GetIntro(),
		Enabled:  in.GetEnabled(),
	}

	err = l.svcCtx.DB.UserDao.Update(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.EditUserOutput{}, nil
}

func (l *EditUserLogic) validate(in *pb.EditUserInput) (*model.User, error) {
	exist, err := l.svcCtx.DB.UserDao.FindOne(l.ctx, in.GetUserId())
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, errUserNotFound
		}
		return nil, err
	}

	maybe, err := l.svcCtx.DB.UserDao.FindOneBySysTypeNickname(l.ctx, exist.SysType, in.GetNickname())
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return nil, err
	}

	if maybe != nil {
		return nil, status.Errorf(codes.AlreadyExists, "用户%s已存在", in.GetNickname())
	}

	return exist, nil
}
