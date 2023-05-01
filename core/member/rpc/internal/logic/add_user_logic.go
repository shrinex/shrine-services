package logic

import (
	"context"
	"core/member/proto/model"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shrine/std/globals"

	"core/member/rpc/internal/svc"
	"core/member/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserLogic {
	return &AddUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddUser 添加用户
func (l *AddUserLogic) AddUser(in *pb.AddUserInput) (*pb.AddUserOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	err = l.validate(in)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		UserId:   l.svcCtx.Leaf.MustNextID(),
		ShopId:   in.GetShopId(),
		SysType:  in.GetSysType(),
		Nickname: in.GetNickname(),
		Avatar:   "https://example.com",
		Intro:    "这个人很懒，什么都没有留下",
		Active:   globals.StatusInactive,
		Enabled:  in.GetEnabled(),
	}
	_, err = l.svcCtx.DB.UserDao.Insert(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.AddUserOutput{
		UserId: user.UserId,
	}, nil
}

func (l *AddUserLogic) validate(in *pb.AddUserInput) error {
	maybe, err := l.svcCtx.DB.UserDao.FindOneBySysTypeNickname(l.ctx, in.GetSysType(), in.GetNickname())
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return err
	}

	if maybe != nil {
		return status.Errorf(codes.AlreadyExists, "用户%s已存在", in.GetNickname())
	}

	return nil
}
