package shop

import (
	"context"
	authcpb "core/authc/rpc/pb"
	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/google/uuid"
	shoppb "unit/shop/rpc/pb"

	"biz/platform/api/internal/svc"
	"biz/platform/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateShopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateShopLogic {
	return &CreateShopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateShopLogic) CreateShop(req *types.CreateShopReq) (resp *types.CreateShopResp, err error) {
	resp = new(types.CreateShopResp)
	shopServer, err := l.svcCtx.Config.ShopRpc.BuildTarget()
	authcServer, err := l.svcCtx.Config.AuthcRpc.BuildTarget()
	if err != nil {
		return
	}

	shop := new(shoppb.AddShopOutput)
	err = dtmgrpc.TccGlobalTransaction(l.svcCtx.Config.Dtm.GRPCServer, uuid.NewString(), func(tcc *dtmgrpc.TccGrpc) error {
		err = tcc.CallBranch(&shoppb.AddShopInput{
			Name:   req.Shop.Name,
			Intro:  req.Shop.Intro,
			Logo:   req.Shop.Logo,
			Status: req.Shop.Status,
			Type:   req.Shop.Type,
		},
			shopServer+shoppb.Service_AddShop_FullMethodName,
			shopServer+shoppb.Service_AddShopConfirm_FullMethodName,
			shopServer+shoppb.Service_AddShopCancel_FullMethodName,
			shop)
		if err != nil {
			return err
		}

		useless := new(authcpb.AddAdminAccountOutput)
		return tcc.CallBranch(&authcpb.AddAdminAccountInput{
			Username: req.Admin.Username,
			Password: req.Admin.Password,
			ShopId:   shop.GetShopId(),
		},
			authcServer+authcpb.Service_AddAdminAccount_FullMethodName,
			authcServer+authcpb.Service_AddAdminAccountConfirm_FullMethodName,
			authcServer+authcpb.Service_AddAdminAccountCancel_FullMethodName,
			useless)
	})

	return
}
