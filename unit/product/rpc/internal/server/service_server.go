// Code generated by goctl. DO NOT EDIT.
// Source: main.proto

package server

import (
	"context"

	"unit/product/rpc/internal/logic"
	"unit/product/rpc/internal/svc"
	"unit/product/rpc/pb"
)

type ServiceServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedServiceServer
}

func NewServiceServer(svcCtx *svc.ServiceContext) *ServiceServer {
	return &ServiceServer{
		svcCtx: svcCtx,
	}
}

// ListCategories 查询所有分类
func (s *ServiceServer) ListCategories(ctx context.Context, in *pb.ListCategoriesInput) (*pb.ListCategoriesOutput, error) {
	l := logic.NewListCategoriesLogic(ctx, s.svcCtx)
	return l.ListCategories(in)
}

// AddCategory 添加分类
func (s *ServiceServer) AddCategory(ctx context.Context, in *pb.AddCategoryInput) (*pb.AddCategoryOutput, error) {
	l := logic.NewAddCategoryLogic(ctx, s.svcCtx)
	return l.AddCategory(in)
}

// EditCategory 编辑分类
func (s *ServiceServer) EditCategory(ctx context.Context, in *pb.EditCategoryInput) (*pb.EditCategoryOutput, error) {
	l := logic.NewEditCategoryLogic(ctx, s.svcCtx)
	return l.EditCategory(in)
}

// RemoveCategory 删除分类
func (s *ServiceServer) RemoveCategory(ctx context.Context, in *pb.RemoveCategoryInput) (*pb.RemoveCategoryOutput, error) {
	l := logic.NewRemoveCategoryLogic(ctx, s.svcCtx)
	return l.RemoveCategory(in)
}
