package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shrine/std/globals"

	"unit/product/rpc/internal/svc"
	"unit/product/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveCategoryLogic {
	return &RemoveCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RemoveCategory 删除分类
func (l *RemoveCategoryLogic) RemoveCategory(in *pb.RemoveCategoryInput) (*pb.RemoveCategoryOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.DB.CategoryDao.FindOne(l.ctx, in.GetCategoryId())
	if errors.Is(err, sqlx.ErrNotFound) {
		return nil, errCategoryNotFound
	}

	if err != nil {
		return nil, err
	}

	err = l.svcCtx.DB.CategoryDao.UpdateStatus(l.ctx, in.GetCategoryId(), globals.StatusRemoved)

	if err != nil {
		return nil, err
	}

	return &pb.RemoveCategoryOutput{}, nil
}
