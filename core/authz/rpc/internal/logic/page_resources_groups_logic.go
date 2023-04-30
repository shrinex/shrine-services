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

type PageResourcesGroupsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageResourcesGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageResourcesGroupsLogic {
	return &PageResourcesGroupsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PageResourcesGroups 分页查询系统中的资源分组
func (l *PageResourcesGroupsLogic) PageResourcesGroups(in *pb.PageResourcesGroupsInput) (*pb.PageResourcesGroupsOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	total, err := l.svcCtx.DB.ResourceGroupDao.CountResourceGroups(l.ctx, in.GetSysType())
	if err != nil {
		return nil, err
	}

	resourceGroups, err := l.svcCtx.DB.ResourceGroupDao.PageResourceGroups(l.ctx, in.GetPageNo(), in.GetPageSize(), in.GetSysType())
	if err != nil {
		return nil, err
	}

	var rows []*pb.ResourceGroup
	_ = copier.CopyWithOption(&rows, resourceGroups, copier.Option{
		Converters: []copier.TypeConverter{conv.NewTime2int64()},
	})

	return &pb.PageResourcesGroupsOutput{
		Rows:  rows,
		Total: total,
		Pages: page.NumPages(total, in.GetPageSize()),
	}, nil
}
