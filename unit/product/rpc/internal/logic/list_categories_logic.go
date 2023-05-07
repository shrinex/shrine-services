package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"shrine/std/utils/slices"
	"unit/product/proto/model"

	"unit/product/rpc/internal/svc"
	"unit/product/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCategoriesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCategoriesLogic {
	return &ListCategoriesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListCategories 查询所有分类
func (l *ListCategoriesLogic) ListCategories(in *pb.ListCategoriesInput) (*pb.ListCategoriesOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	categories, err := l.svcCtx.DB.CategoryDao.ListCategories(l.ctx, in.GetName())
	if err != nil {
		return nil, err
	}

	return &pb.ListCategoriesOutput{
		Categories: l.treeify(categories),
	}, nil
}

func (l *ListCategoriesLogic) treeify(categories []*model.Category) []*pb.CategoryNode {
	nodes := slices.Map(categories, func(e *model.Category) (ret *pb.CategoryNode) {
		ret = new(pb.CategoryNode)
		_ = copier.Copy(ret, e)
		return
	})

	lookup := slices.GroupingBy(nodes, func(e *pb.CategoryNode) int64 {
		return e.ParentId
	})

	for _, e := range nodes {
		if children, ok := lookup[e.GetCategoryId()]; ok {
			e.Children = children
		} else {
			e.Children = []*pb.CategoryNode{}
		}
	}

	return slices.Filter(nodes, func(e *pb.CategoryNode) bool {
		return e.GetParentId() == rootCategoryId
	})
}
