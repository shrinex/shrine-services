package logic

import (
	"context"
	"core/authz/proto/model"
	"core/authz/rpc/internal/svc"
	"core/authz/rpc/pb"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddResourceLogic {
	return &AddResourceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddResource 添加资源
func (l *AddResourceLogic) AddResource(in *pb.AddResourceInput) (*pb.AddResourceOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	err = l.validate(in)
	if err != nil {
		return nil, err
	}

	// 添加资源
	resource := &model.Resource{
		ResourceId: l.svcCtx.Leaf.MustNextID(),
		GroupId:    in.GetGroupId(),
		Name:       in.GetName(),
		Method:     in.GetMethod(),
		Pattern:    in.GetPattern(),
		SysType:    in.GetSysType(),
	}
	_, err = l.svcCtx.DB.ResourceDao.Insert(l.ctx, resource)
	if err != nil {
		return nil, err
	}

	// 清除管理员资源缓存
	err = l.svcCtx.Cache.ResourceCache.ClearResourcesBySysType(l.ctx, in.GetSysType())
	if err != nil {
		return nil, err
	}

	return &pb.AddResourceOutput{
		ResourceId: resource.ResourceId,
	}, nil
}

func (l *AddResourceLogic) validate(in *pb.AddResourceInput) error {
	maybe, err := l.svcCtx.DB.ResourceDao.FindOneByName(l.ctx, in.GetName())
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return err
	}

	if maybe != nil {
		return status.Errorf(codes.AlreadyExists, "分组%s已存在", in.GetName())
	}

	maybe, err = l.svcCtx.DB.ResourceDao.FindOneBySysTypeMethodPattern(l.ctx,
		in.GetSysType(), in.GetMethod(), in.GetPattern())
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return err
	}

	if maybe != nil {
		return status.Errorf(codes.AlreadyExists, "分组%s已存在", in.GetName())
	}

	return nil
}
