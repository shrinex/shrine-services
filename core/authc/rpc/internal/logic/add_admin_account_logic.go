package logic

import (
	"context"
	"core/authc/proto/model"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/crypto/bcrypt"
	"shrine/std/globals"
	"shrine/std/utils/dtmx"
	"shrine/std/utils/sqle"

	"core/authc/rpc/internal/svc"
	"core/authc/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAdminAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAdminAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAdminAccountLogic {
	return &AddAdminAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddAdminAccount 添加商家端管理员用户
func (l *AddAdminAccountLogic) AddAdminAccount(in *pb.AddAdminAccountInput) (*pb.AddAdminAccountOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, dtmx.Abort(err)
	}

	err = l.validate(in)
	if err != nil {
		return nil, dtmx.Abort(err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, dtmx.Abort(err)
	}

	userId := l.svcCtx.Leaf.MustNextID()
	accountId := l.svcCtx.Leaf.MustNextID()
	barrier := dtmx.MustBarrierFromGrpc(l.ctx)
	err = barrier.CallWithDB(l.svcCtx.DB.RawDB(), func(tx *sql.Tx) error {
		txSession := sqlx.NewSessionFromTx(tx)
		_, err = l.svcCtx.DB.UserDao.TxInsert(l.ctx, txSession, &model.User{
			UserId:   userId,
			ShopId:   in.GetShopId(),
			SysType:  globals.SysTypeMerchant,
			Nickname: in.GetUsername(),
			Avatar:   "https://example.com",
			Intro:    "这个人很懒，什么都没有留下",
			Enabled:  globals.FlagTrue,
		})
		if err != nil {
			if sqle.Is(err, sqle.DuplicateEntry) {
				return dtmx.Abortf(errUserExistsDesc)
			}
			return err
		}

		_, err = l.svcCtx.DB.AccountDao.TxInsert(l.ctx, txSession, &model.Account{
			AccountId: accountId,
			UserId:    userId,
			Username:  in.GetUsername(),
			Password:  string(password),
			SysType:   globals.SysTypeMerchant,
			ShopId:    in.GetShopId(),
			IsAdmin:   globals.FlagTrue,
			Enabled:   globals.FlagTrue,
		})
		if err != nil {
			if sqle.Is(err, sqle.DuplicateEntry) {
				return dtmx.Abortf(errAccountExistsDesc)
			}
			return err
		}

		return nil
	})

	if err != nil {
		return nil, dtmx.Retry(err)
	}

	return &pb.AddAdminAccountOutput{
		UserId:    userId,
		AccountId: accountId,
	}, nil
}

func (l *AddAdminAccountLogic) validate(in *pb.AddAdminAccountInput) error {
	exists, err := l.svcCtx.DB.AccountDao.AccountExistsBySysTypeAndUsername(l.ctx, globals.SysTypeMerchant, in.GetUsername())
	if err != nil {
		return err
	}

	if exists {
		return errAccountExists
	}

	exists, err = l.svcCtx.DB.UserDao.UserExistsBySysTypeAndNickname(l.ctx, globals.SysTypeMerchant, in.GetUsername())
	if err != nil {
		return err
	}

	if exists {
		return errUserExists
	}

	return nil
}
