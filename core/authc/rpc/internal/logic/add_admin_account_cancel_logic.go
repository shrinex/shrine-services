package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shrine/std/globals"
	"shrine/std/utils/dtmx"

	"core/authc/rpc/internal/svc"
	"core/authc/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAdminAccountCancelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAdminAccountCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAdminAccountCancelLogic {
	return &AddAdminAccountCancelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddAdminAccountCancel 添加商家端管理员用户回滚
func (l *AddAdminAccountCancelLogic) AddAdminAccountCancel(in *pb.AddAdminAccountInput) (*pb.AddAdminAccountOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	logx.Info("calling add admin account revert...")
	barrier := dtmx.MustBarrierFromGrpc(l.ctx)
	err = barrier.CallWithDB(l.svcCtx.DB.RawDB(), func(tx *sql.Tx) error {
		txSession := sqlx.NewSessionFromTx(tx)
		account, rerr := l.svcCtx.DB.AccountDao.TxFindOneBySysTypeUsername(l.ctx, txSession, globals.SysTypeMerchant, in.GetUsername())
		if errors.Is(rerr, sqlx.ErrNotFound) {
			return nil
		}

		if rerr != nil {
			return rerr
		}

		rerr = l.svcCtx.DB.AccountDao.TxDelete(l.ctx, txSession, account.AccountId)
		if rerr != nil {
			return rerr
		}

		user, rerr := l.svcCtx.DB.UserDao.TxFindOneBySysTypeNickname(l.ctx, txSession, globals.SysTypeMerchant, in.GetUsername())
		if errors.Is(rerr, sqlx.ErrNotFound) {
			return nil
		}

		if rerr != nil {
			return rerr
		}

		rerr = l.svcCtx.DB.UserDao.TxDelete(l.ctx, txSession, user.UserId)
		if rerr != nil {
			return rerr
		}

		return nil
	})

	if err != nil {
		return nil, dtmx.Retry(err)
	}

	return &pb.AddAdminAccountOutput{}, nil
}
