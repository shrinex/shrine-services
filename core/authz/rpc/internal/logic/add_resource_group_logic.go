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

type AddResourceGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddResourceGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddResourceGroupLogic {
	return &AddResourceGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddResourceGroup 添加资源分组
func (l *AddResourceGroupLogic) AddResourceGroup(in *pb.AddResourceGroupInput) (*pb.AddResourceGroupOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	err = l.validate(in)
	if err != nil {
		return nil, err
	}

	group := &model.ResourceGroup{
		GroupId: l.svcCtx.Leaf.MustNextID(),
		Name:    in.GetName(),
		Remark:  in.GetRemark(),
		SysType: in.GetSysType(),
	}
	_, err = l.svcCtx.DB.ResourceGroupDao.Insert(l.ctx, group)
	if err != nil {
		return nil, err
	}

	return &pb.AddResourceGroupOutput{
		GroupId: group.GroupId,
	}, nil
}

func (l *AddResourceGroupLogic) validate(in *pb.AddResourceGroupInput) error {
	maybe, err := l.svcCtx.DB.ResourceGroupDao.FindOneByName(l.ctx, in.GetName())
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return err
	}

	if maybe != nil {
		return status.Errorf(codes.AlreadyExists, "分组%s已存在", in.GetName())
	}

	return nil
}
