package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"shrine/std/conv"
	"shrine/std/utils/page"

	"core/authz/rpc/internal/svc"
	"core/authz/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageMenusLogic {
	return &PageMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PageMenus 分页查询系统中的菜单列表
func (l *PageMenusLogic) PageMenus(in *pb.PageMenusInput) (*pb.PageMenusOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	total, err := l.svcCtx.DB.MenuDao.CountMenus(l.ctx, in.GetSysType())
	if err != nil {
		return nil, err
	}

	menus, err := l.svcCtx.DB.MenuDao.PageMenus(l.ctx, in.GetPageNo(), in.GetPageSize(), in.GetSysType())
	if err != nil {
		return nil, err
	}

	var rows []*pb.Menu
	_ = copier.CopyWithOption(&rows, menus, copier.Option{
		Converters: []copier.TypeConverter{conv.NewTime2int64()},
	})

	return &pb.PageMenusOutput{
		Rows:  rows,
		Total: total,
		Pages: page.NumPages(total, in.GetPageSize()),
	}, nil
}
