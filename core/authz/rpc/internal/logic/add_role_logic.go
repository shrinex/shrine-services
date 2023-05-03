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

type AddRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRoleLogic {
	return &AddRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddRole 添加角色
func (l *AddRoleLogic) AddRole(in *pb.AddRoleInput) (*pb.AddRoleOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	err = l.validate(in)
	if err != nil {
		return nil, err
	}

	// 添加角色
	role := &model.Role{
		RoleId:    l.svcCtx.Leaf.MustNextID(),
		CreatorId: in.GetCreatorId(),
		Name:      in.GetName(),
		Remark:    in.GetRemark(),
		ShopId:    in.GetShopId(),
		SysType:   in.GetSysType(),
	}
	_, err = l.svcCtx.DB.RoleDao.Insert(l.ctx, role)
	if err != nil {
		return nil, err
	}

	// 清除管理员角色缓存
	err = l.svcCtx.Cache.RoleCache.ClearRolesBySysType(l.ctx, in.GetSysType())
	if err != nil {
		return nil, err
	}

	return &pb.AddRoleOutput{
		RoleId: role.RoleId,
	}, nil
}

func (l *AddRoleLogic) validate(in *pb.AddRoleInput) error {
	maybe, err := l.svcCtx.DB.RoleDao.FindOneByName(l.ctx, in.GetName())
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return err
	}

	if maybe != nil {
		return status.Errorf(codes.AlreadyExists, "角色%s已存在", in.GetName())
	}

	return nil
}
