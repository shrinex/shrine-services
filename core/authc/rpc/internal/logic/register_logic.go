package logic

import (
	"context"
	"core/authc/proto/model"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shrine/std/globals"
	"shrine/std/utils/dtmx"
	"shrine/std/utils/verify"

	"core/authc/rpc/internal/svc"
	"core/authc/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Register 用户注册
func (l *RegisterLogic) Register(in *pb.RegisterInput) (*pb.RegisterOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	err = l.validate(in)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	password, err := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	userId := l.svcCtx.Leaf.MustNextID()
	accountId := l.svcCtx.Leaf.MustNextID()
	barrier := dtmx.MustBarrierFromGrpc(l.ctx)
	err = barrier.CallWithDB(l.svcCtx.DB.RawDB, func(tx *sql.Tx) error {
		_, err = l.svcCtx.DB.UserDao.TxInsert(l.ctx, tx, &model.User{
			UserId:   l.svcCtx.Leaf.MustNextID(),
			ShopId:   in.GetShopId(),
			SysType:  in.GetSysType(),
			Nickname: in.GetNickname(),
			Avatar:   "https://example.com",
			Intro:    "这个人很懒，什么都没有留下",
			Active:   globals.StatusInactive,
			Enabled:  globals.FlagTrue,
		})
		if err != nil {
			if verify.Duplicated(err) {
				return status.Errorf(codes.Aborted, errUserExistsDesc)
			}
			return err
		}

		_, err = l.svcCtx.DB.AccountDao.TxInsert(l.ctx, tx, &model.Account{
			AccountId: l.svcCtx.Leaf.MustNextID(),
			UserId:    userId,
			Username:  in.GetUsername(),
			Password:  string(password),
			SysType:   in.GetSysType(),
			ShopId:    in.GetShopId(),
			IsAdmin:   in.GetIsAdmin(),
			Enabled:   globals.FlagTrue,
		})
		if err != nil {
			if verify.Duplicated(err) {
				return status.Errorf(codes.Aborted, errAccountExistsDesc)
			}
			return err
		}

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterOutput{
		UserId:    userId,
		AccountId: accountId,
	}, nil
}

func (l *RegisterLogic) validate(in *pb.RegisterInput) error {
	exists, err := l.svcCtx.DB.AccountDao.AccountExistsBySysTypeAndUsername(l.ctx, in.GetSysType(), in.GetUsername())
	if err != nil {
		return err
	}

	if exists {
		return errAccountExists
	}

	exists, err = l.svcCtx.DB.UserDao.UserExistsBySysTypeAndNickname(l.ctx, in.GetSysType(), in.GetNickname())
	if err != nil {
		return err
	}

	if exists {
		return errUserExists
	}

	return nil
}
