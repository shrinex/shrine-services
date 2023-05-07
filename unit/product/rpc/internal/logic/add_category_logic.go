package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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

	level := int64(1)
	if in.GetParentId() != rootCategoryId {
		parent, err := l.svcCtx.DB.CategoryDao.FindOne(l.ctx, in.GetParentId())
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, errParentCategoryNotFound
		}

		if err != nil {
			return nil, err
		}

		if parent.Level >= maxCategoryLevel {
			return nil, errCategoryLevelOverflow
		}

		level = parent.Level + 1
	}

	catId := l.svcCtx.Leaf.MustNextID()
	_, err = l.svcCtx.DB.CategoryDao.Insert(l.ctx, &model.Category{
		CategoryId: catId,
		ParentId:   in.GetParentId(),
		Name:       in.GetName(),
		Remark:     in.GetRemark(),
		Icon:       in.GetIcon(),
		Level:      level,
		Status:     in.GetStatus(),
		Weight:     in.GetWeight(),
	})

	if err != nil {
		return nil, err
	}

	return &pb.AddCategoryOutput{
		CategoryId: catId,
	}, nil
}
