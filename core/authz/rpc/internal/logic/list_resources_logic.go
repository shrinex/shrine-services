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

type ListResourcesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListResourcesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListResourcesLogic {
	return &ListResourcesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListResources 查询用户拥有的资源列表
func (l *ListResourcesLogic) ListResources(in *pb.ListResourcesInput) (*pb.ListResourcesOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	resources, err := l.svcCtx.Cache.ResourceCache.ListResources(l.ctx, in.GetSysType(), in.GetUserId(), in.GetIsAdmin())
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		return nil, err
	}

	var outputs []*pb.Resource
	_ = copier.CopyWithOption(&outputs, resources, copier.Option{
		Converters: []copier.TypeConverter{conv.NewTime2int64()},
	})

	return &pb.ListResourcesOutput{
		Resources: outputs,
	}, nil
}
