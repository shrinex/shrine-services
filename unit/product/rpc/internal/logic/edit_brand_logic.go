package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shrine/std/utils/slices"
	"unit/product/proto/model"
	"unit/product/rpc/internal/svc"
	"unit/product/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditBrandLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditBrandLogic {
	return &EditBrandLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// EditBrand 编辑品牌
func (l *EditBrandLogic) EditBrand(in *pb.EditBrandInput) (*pb.EditBrandOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	brand, err := l.svcCtx.DB.BrandDao.FindOne(l.ctx, in.GetBrandId())
	if errors.Is(err, sqlx.ErrNotFound) {
		return nil, errBrandNotFound
	}

	if err != nil {
		return nil, err
	}

	err = l.svcCtx.DB.RawConn.TransactCtx(l.ctx, func(ctx context.Context, tx sqlx.Session) error {
		err = l.svcCtx.DB.BrandDao.TxUpdate(l.ctx, tx, &model.Brand{
			BrandId: brand.BrandId,
			Name:    in.GetName(),
			Remark:  in.GetRemark(),
			Logo:    in.GetLogo(),
			Status:  in.GetStatus(),
			Weight:  in.GetWeight(),
		})

		if err != nil {
			return err
		}

		err = l.svcCtx.DB.BrandCategoryDao.TxDeleteByBrandId(l.ctx, tx, brand.BrandId)
		if err != nil {
			return err
		}

		_, err = l.svcCtx.DB.BrandCategoryDao.TxInsertBatch(l.ctx, tx,
			slices.Map(in.GetCategoryIds(),
				func(e int64) *model.BrandCategoryRel {
					return &model.BrandCategoryRel{
						RelId:      l.svcCtx.Leaf.MustNextID(),
						BrandId:    brand.BrandId,
						CategoryId: e,
					}
				},
			),
		)

		return err
	})

	if err != nil {
		return nil, err
	}

	return &pb.EditBrandOutput{}, nil
}
