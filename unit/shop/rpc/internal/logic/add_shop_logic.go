package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, err
	}

	err = l.validate(in)
	if err != nil {
		return nil, err
	}

	shop := &model.Shop{
		ShopId: l.svcCtx.Leaf.MustNextID(),
		Name:   in.GetName(),
		Intro:  in.GetIntro(),
		Logo:   in.GetLogo(),
		Status: in.GetStatus(),
		Type:   in.GetType(),
	}
	_, err = l.svcCtx.DB.ShopDao.Insert(l.ctx, shop)
	if err != nil {
		return nil, err
	}

	return &pb.AddShopOutput{
		ShopId: shop.ShopId,
	}, nil
}

func (l *AddShopLogic) validate(in *pb.AddShopInput) error {
	maybe, err := l.svcCtx.DB.ShopDao.FindOneByName(l.ctx, in.GetName())
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return err
	}

	if maybe != nil {
		return status.Errorf(codes.AlreadyExists, "店铺%s已存在", in.GetName())
	}

	return nil
}
