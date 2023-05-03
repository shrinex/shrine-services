package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shrine/std/utils/dtmx"
	"shrine/std/utils/verify"
	"unit/shop/proto/model"

	"unit/shop/rpc/internal/svc"
	"unit/shop/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddShopLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddShopLogic {
	return &AddShopLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddShop 创建店铺
func (l *AddShopLogic) AddShop(in *pb.AddShopInput) (*pb.AddShopOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	err = l.validate(in)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	shopId := l.svcCtx.Leaf.MustNextID()
	barrier := dtmx.MustBarrierFromGrpc(l.ctx)
	err = barrier.CallWithDB(l.svcCtx.DB.RawDB, func(tx *sql.Tx) error {
		_, err = l.svcCtx.DB.ShopDao.TxInsert(l.ctx, tx, &model.Shop{
			ShopId: shopId,
			Name:   in.GetName(),
			Intro:  in.GetIntro(),
			Logo:   in.GetLogo(),
			Status: in.GetStatus(),
			Type:   in.GetType(),
		})
		if err != nil {
			if verify.Duplicated(err) {
				return status.Errorf(codes.Aborted, errShopExistsDesc)
			}
			return err
		}

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AddShopOutput{
		ShopId: shopId,
	}, nil
}

func (l *AddShopLogic) validate(in *pb.AddShopInput) error {
	maybe, err := l.svcCtx.DB.ShopDao.FindOneByName(l.ctx, in.GetName())
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return err
	}

	if maybe != nil {
		return errShopExists
	}

	return nil
}
