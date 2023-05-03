package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shrine/std/utils/dtmx"

	"core/authc/rpc/internal/svc"
	"core/authc/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterCancelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterCancelLogic {
	return &RegisterCancelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RegisterCancel 用户注册回滚
func (l *RegisterCancelLogic) RegisterCancel(in *pb.RegisterInput) (*pb.RegisterOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	logx.Info("calling register revert...")
	barrier := dtmx.MustBarrierFromGrpc(l.ctx)
	err = barrier.CallWithDB(l.svcCtx.DB.RawDB, func(tx *sql.Tx) error {
		account, rerr := l.svcCtx.DB.AccountDao.TxFindOneBySysTypeUsername(l.ctx, tx, in.GetSysType(), in.GetUsername())
		if errors.Is(rerr, sqlx.ErrNotFound) {
			return nil
		}

		if rerr != nil {
			return rerr
		}

		rerr = l.svcCtx.DB.AccountDao.TxDelete(l.ctx, tx, account.AccountId)
		if rerr != nil {
			return rerr
		}

		user, rerr := l.svcCtx.DB.UserDao.TxFindOneBySysTypeNickname(l.ctx, tx, in.GetSysType(), in.GetUsername())
		if errors.Is(rerr, sqlx.ErrNotFound) {
			return nil
		}

		if rerr != nil {
			return rerr
		}

		rerr = l.svcCtx.DB.UserDao.TxDelete(l.ctx, tx, user.UserId)
		if rerr != nil {
			return rerr
		}

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterOutput{}, nil
}
