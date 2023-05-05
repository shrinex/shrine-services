package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shrine/std/utils/dtmx"

	"unit/shop/rpc/internal/svc"
	"unit/shop/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddShopCancelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddShopCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddShopCancelLogic {
	return &AddShopCancelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddShopCancel 创建店铺回滚
func (l *AddShopCancelLogic) AddShopCancel(in *pb.AddShopInput) (*pb.AddShopOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	logx.Info("calling add shop revert...")
	barrier := dtmx.MustBarrierFromGrpc(l.ctx)
	err = barrier.CallWithDB(l.svcCtx.DB.RawDB, func(tx *sql.Tx) error {
		txSession := sqlx.NewSessionFromTx(tx)
		shop, rerr := l.svcCtx.DB.ShopDao.TxFindOneByName(l.ctx, txSession, in.GetName())
		if errors.Is(rerr, sqlx.ErrNotFound) {
			return nil
		}

		if rerr != nil {
			return rerr
		}

		return l.svcCtx.DB.ShopDao.TxDelete(l.ctx, txSession, shop.ShopId)
	})

	if err != nil {
		return nil, dtmx.Retry(err)
	}

	return &pb.AddShopOutput{}, nil
}
