package logic

import (
	"context"
	"core/authc/proto/model"
	"core/authc/rpc/internal/svc"
	"core/authc/rpc/pb"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/crypto/bcrypt"
	"shrine/std/globals"
	"shrine/std/utils/sqle"

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
		return nil, err
	}

	err = l.validate(in)
	if err != nil {
		return nil, err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userId := l.svcCtx.Leaf.MustNextID()
	accountId := l.svcCtx.Leaf.MustNextID()
	err = l.svcCtx.DB.RawConn.TransactCtx(l.ctx, func(ctx context.Context, tx sqlx.Session) error {
		_, err = l.svcCtx.DB.UserDao.TxInsert(l.ctx, tx, &model.User{
			UserId:   userId,
			ShopId:   in.GetShopId(),
			SysType:  in.GetSysType(),
			Nickname: in.GetUsername(),
			Avatar:   "https://example.com",
			Intro:    "这个人很懒，什么都没有留下",
			Enabled:  globals.FlagTrue,
		})
		if err != nil {
			if sqle.Is(err, sqle.DuplicateEntry) {
				return errUserExists
			}
			return err
		}

		_, err = l.svcCtx.DB.AccountDao.TxInsert(l.ctx, tx, &model.Account{
			AccountId: accountId,
			UserId:    userId,
			Username:  in.GetUsername(),
			Password:  string(password),
			SysType:   in.GetSysType(),
			ShopId:    in.GetShopId(),
			IsAdmin:   globals.FlagFalse,
			Enabled:   globals.FlagTrue,
		})
		if err != nil {
			if sqle.Is(err, sqle.DuplicateEntry) {
				return errAccountExists
			}
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
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

	exists, err = l.svcCtx.DB.UserDao.UserExistsBySysTypeAndNickname(l.ctx, in.GetSysType(), in.GetUsername())
	if err != nil {
		return err
	}

	if exists {
		return errUserExists
	}

	return nil
}
