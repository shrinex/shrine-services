package logic

import (
	"context"
	"core/authz/proto/model"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/mathx"
	"github.com/zeromicro/go-zero/core/mr"
	"shrine/std/utils/iter"
	"shrine/std/utils/page"

	"core/authz/rpc/internal/svc"
	"core/authz/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveResourceLogic {
	return &RemoveResourceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RemoveResource 删除资源
func (l *RemoveResourceLogic) RemoveResource(in *pb.RemoveResourceInput) (*pb.RemoveResourceOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	resource, err := l.svcCtx.DB.ResourceDao.FindOne(l.ctx, in.GetResourceId())
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return &pb.RemoveResourceOutput{}, nil
		}
		return nil, err
	}

	// 删除 角色->资源 缓存
	var perr error
	pager := iter.NewPager(func() (int64, error) {
		return l.svcCtx.DB.RoleResourceDao.CountRoleIdsByResourceId(l.ctx, in.GetResourceId())
	}, func(offset int64, size int64) ([]int64, error) {
		return l.svcCtx.DB.RoleResourceDao.PageRoleIdsByResourceId(l.ctx, offset, size, in.GetResourceId())
	})
	err = mr.MapReduceVoid(func(source chan<- []int64) {
		for {
			var items []int64
			items, perr = pager.Next()
			if errors.Is(perr, iter.Done) {
				perr = nil // do not treat iter.Done as error
				break
			}

			if perr != nil {
				l.Logger.Errorf("查询资源关联的角色ID失败, resourceId: %d, err: %v", in.GetResourceId(), perr)
				return
			}

			source <- items
		}
	}, func(items []int64, writer mr.Writer[[]int64], cancel func(error)) {
		writer.Write(items)
	}, func(pipe <-chan []int64, cancel func(error)) {
		for items := range pipe {
			rerr := l.svcCtx.Cache.ResourceCache.ClearResources(l.ctx, resource.SysType, items...)
			if rerr != nil {
				l.Logger.Errorf("清除用户角色缓存失败: %v", rerr)
				return
			}
		}
	}, mr.WithWorkers(mathx.MinInt(16, int(page.NumPages(pager.EstimatedSize(), 256)))))

	if perr != nil {
		return nil, perr
	}

	if err != nil {
		return nil, err
	}

	// 删除资源
	err = l.svcCtx.DB.ResourceDao.Delete(l.ctx, in.GetResourceId())
	if err != nil {
		return nil, err
	}

	return &pb.RemoveResourceOutput{}, nil
}
