package logic

import (
	"context"
	"core/authz/rpc/internal/repo/cache"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"shrine/std/conv"

	"core/authz/rpc/internal/svc"
	"core/authz/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenusLogic {
	return &ListMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListMenus 查询用户拥有的菜单列表
func (l *ListMenusLogic) ListMenus(in *pb.ListMenusInput) (*pb.ListMenusOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	roles, err := l.svcCtx.Cache.MenuCache.ListMenus(l.ctx, in.GetSysType(), in.GetIsAdmin(), in.GetUserId())
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		return nil, err
	}

	var outputs []*pb.Menu
	_ = copier.CopyWithOption(&outputs, roles, copier.Option{
		Converters: []copier.TypeConverter{conv.NewTime2int64()},
	})

	return &pb.ListMenusOutput{
		Menus: outputs,
	}, nil
}
