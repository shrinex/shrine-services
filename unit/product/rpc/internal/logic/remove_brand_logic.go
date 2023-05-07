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

type RemoveBrandLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveBrandLogic {
	return &RemoveBrandLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RemoveBrand 删除品牌
func (l *RemoveBrandLogic) RemoveBrand(in *pb.RemoveBrandInput) (*pb.RemoveBrandOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.DB.BrandDao.FindOne(l.ctx, in.GetBrandId())
	if errors.Is(err, sqlx.ErrNotFound) {
		return nil, errBrandNotFound
	}

	if err != nil {
		return nil, err
	}

	err = l.svcCtx.DB.BrandDao.UpdateStatus(l.ctx, in.GetBrandId(), globals.StatusRemoved)

	if err != nil {
		return nil, err
	}

	return &pb.RemoveBrandOutput{}, nil
}
