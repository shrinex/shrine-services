package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		shop, rerr := l.svcCtx.DB.ShopDao.TxFindOneByName(l.ctx, tx, in.GetName())
		if errors.Is(rerr, sqlx.ErrNotFound) {
			return nil
		}

		if rerr != nil {
			return rerr
		}

		return l.svcCtx.DB.ShopDao.TxDelete(l.ctx, tx, shop.ShopId)
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AddShopOutput{}, nil
}
