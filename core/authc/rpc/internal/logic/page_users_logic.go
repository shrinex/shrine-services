package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"shrine/std/conv"
	"shrine/std/utils/page"
	"shrine/std/utils/slices"

	"core/authc/rpc/internal/svc"
	"core/authc/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageUsersLogic {
	return &PageUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PageUsers 分页获取用户列表
func (l *PageUsersLogic) PageUsers(in *pb.PageUsersInput) (*pb.PageUsersOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	total, err := l.svcCtx.DB.UserDao.CountUsers(l.ctx, in.GetSysType(), in.GetShopId(), in.GetNickname())
	if err != nil {
		return nil, err
	}

	if total <= 0 {
		return &pb.PageUsersOutput{
			Pages: 0,
			Total: 0,
			Rows:  slices.Empty[*pb.User](),
		}, nil
	}

	users, err := l.svcCtx.DB.UserDao.PageUsers(l.ctx, in.GetPageNo(), in.GetPageSize(),
		in.GetSysType(), in.GetShopId(), in.GetNickname())
	if err != nil {
		return nil, err
	}

	var rows []*pb.User
	_ = copier.CopyWithOption(&rows, users, copier.Option{
		Converters: []copier.TypeConverter{conv.NewTime2int64()},
	})

	return &pb.PageUsersOutput{
		Rows:  rows,
		Total: total,
		Pages: page.NumPages(total, in.GetPageSize()),
	}, nil
}
