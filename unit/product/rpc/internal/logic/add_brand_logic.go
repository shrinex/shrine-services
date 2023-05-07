package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shrine/std/utils/slices"
	"shrine/std/utils/sqle"
	"unit/product/proto/model"

	"unit/product/rpc/internal/svc"
	"unit/product/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddBrandLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBrandLogic {
	return &AddBrandLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddBrand 添加品牌
func (l *AddBrandLogic) AddBrand(in *pb.AddBrandInput) (*pb.AddBrandOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	err = l.validate(in)
	if err != nil {
		return nil, err
	}

	brandId := l.svcCtx.Leaf.MustNextID()
	err = l.svcCtx.DB.RawConn.TransactCtx(l.ctx, func(ctx context.Context, tx sqlx.Session) error {
		_, err = l.svcCtx.DB.BrandDao.TxInsert(l.ctx, tx, &model.Brand{
			BrandId: brandId,
			Name:    in.GetName(),
			Remark:  in.GetRemark(),
			Logo:    in.GetLogo(),
			Status:  in.GetStatus(),
			Weight:  in.GetWeight(),
		})

		if err != nil {
			if sqle.Is(err, sqle.DuplicateEntry) {
				return errBrandExists
			}
			return err
		}

		_, err = l.svcCtx.DB.BrandCategoryDao.TxInsertBatch(l.ctx, tx,
			slices.Map(in.GetCategoryIds(),
				func(e int64) *model.BrandCategoryRel {
					return &model.BrandCategoryRel{
						RelId:      l.svcCtx.Leaf.MustNextID(),
						BrandId:    brandId,
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

	return &pb.AddBrandOutput{
		BrandId: brandId,
	}, nil
}

func (l *AddBrandLogic) validate(in *pb.AddBrandInput) error {
	maybe, err := l.svcCtx.DB.BrandDao.FindOneByName(l.ctx, in.GetName())
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return err
	}

	if maybe != nil {
		return errBrandExists
	}

	return nil
}
