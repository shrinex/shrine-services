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

type EditCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditCategoryLogic {
	return &EditCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// EditCategory 编辑分类
func (l *EditCategoryLogic) EditCategory(in *pb.EditCategoryInput) (*pb.EditCategoryOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	category, err := l.svcCtx.DB.CategoryDao.FindOne(l.ctx, in.GetCategoryId())
	if errors.Is(err, sqlx.ErrNotFound) {
		return nil, errCategoryNotFound
	}

	if err != nil {
		return nil, err
	}

	err = l.svcCtx.DB.CategoryDao.Update(l.ctx, &model.Category{
		CategoryId: category.CategoryId,
		ParentId:   category.ParentId,
		Name:       in.GetName(),
		Remark:     in.GetRemark(),
		Icon:       in.GetIcon(),
		Level:      category.Level,
		Status:     in.GetStatus(),
		Weight:     in.GetWeight(),
	})

	if err != nil {
		return nil, err
	}

	return &pb.EditCategoryOutput{}, nil
}
