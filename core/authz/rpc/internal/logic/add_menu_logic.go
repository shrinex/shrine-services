package logic

import (
	"context"
	"core/authz/proto/model"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"core/authz/rpc/internal/svc"
	"core/authz/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMenuLogic {
	return &AddMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddMenu 添加菜单
func (l *AddMenuLogic) AddMenu(in *pb.AddMenuInput) (*pb.AddMenuOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	err = l.validate(in)
	if err != nil {
		return nil, err
	}

	var level int64 = 1
	menuId := l.svcCtx.Leaf.MustNextID()
	path := fmt.Sprintf("/%d", menuId)
	if in.GetParentId() != rootMenuID {
		parent, err := l.svcCtx.DB.MenuDao.FindOne(l.ctx, in.GetParentId())
		if err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				return nil, errParentMenuNotExists
			}
			return nil, err
		}

		level = parent.Level + 1
		path = fmt.Sprintf("%s/%d", parent.Path, menuId)
	}

	_, err = l.svcCtx.DB.MenuDao.Insert(l.ctx, &model.Menu{
		MenuId:   menuId,
		Name:     in.GetName(),
		Icon:     in.GetIcon(),
		ParentId: in.GetParentId(),
		SysType:  in.GetSysType(),
		Weight:   in.GetWeight(),
		Level:    level,
		Path:     path,
	})
	if err != nil {
		return nil, err
	}

	// 清除管理员资源缓存
	err = l.svcCtx.Cache.MenuCache.ClearMenusBySysType(l.ctx, in.GetSysType())
	if err != nil {
		return nil, err
	}

	return &pb.AddMenuOutput{
		MenuId: menuId,
	}, nil
}

func (l *AddMenuLogic) validate(in *pb.AddMenuInput) error {
	maybe, err := l.svcCtx.DB.MenuDao.FindOneByName(l.ctx, in.GetName())
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return err
	}

	if maybe != nil {
		return status.Errorf(codes.AlreadyExists, "菜单%s已存在", in.GetName())
	}

	return nil
}
