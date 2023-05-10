package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shrine/std/utils/sqle"
	"unit/product/proto/model"

	"unit/product/rpc/internal/svc"
	"unit/product/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCategoryLogic {
	return &AddCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddCategory 添加分类
func (l *AddCategoryLogic) AddCategory(in *pb.AddCategoryInput) (*pb.AddCategoryOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	level, err := l.validate(in)
	if err != nil {
		return nil, err
	}

	catId := l.svcCtx.Leaf.MustNextID()
	_, err = l.svcCtx.DB.CategoryDao.Insert(l.ctx, &model.Category{
		CategoryId: catId,
		ParentId:   in.GetParentId(),
		GroupId:    in.GetGroupId(),
		Name:       in.GetName(),
		Remark:     in.GetRemark(),
		Icon:       in.GetIcon(),
		Level:      level,
		Status:     in.GetStatus(),
		Weight:     in.GetWeight(),
	})

	if err != nil {
		if sqle.Is(err, sqle.DuplicateEntry) {
			return nil, errCategoryExists
		}
		return nil, err
	}

	return &pb.AddCategoryOutput{
		CategoryId: catId,
	}, nil
}

func (l *AddCategoryLogic) validate(in *pb.AddCategoryInput) (int64, error) {
	if in.GetParentId() == rootCategoryId {
		return 1, nil
	}

	parent, err := l.svcCtx.DB.CategoryDao.FindOne(l.ctx, in.GetParentId())
	if errors.Is(err, sqlx.ErrNotFound) {
		return 0, errParentCategoryNotFound
	}

	if err != nil {
		return 0, err
	}

	if parent.Level >= maxCategoryLevel {
		return 0, errCategoryLevelOverflow
	}

	return parent.Level + 1, nil
}
