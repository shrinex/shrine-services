package logic

import (
	"context"
	"core/authz/proto/model"
	"core/authz/rpc/internal/svc"
	"core/authz/rpc/pb"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/mathx"
	"github.com/zeromicro/go-zero/core/mr"
	"shrine/std/utils/iter"
	"shrine/std/utils/page"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveRoleLogic {
	return &RemoveRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RemoveRole 删除角色
func (l *RemoveRoleLogic) RemoveRole(in *pb.RemoveRoleInput) (*pb.RemoveRoleOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	role, err := l.svcCtx.DB.RoleDao.FindOne(l.ctx, in.GetRoleId())
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return &pb.RemoveRoleOutput{}, nil
		}
		return nil, err
	}

	// 删除 角色->权限 缓存
	err = l.svcCtx.Cache.ResourceCache.ClearResourcesByRoleIds(l.ctx, role.SysType, role.RoleId)
	if err != nil {
		return nil, err
	}

	// 删除 用户->角色 缓存
	var perr error
	pager := iter.NewPager(func() (int64, error) {
		return l.svcCtx.DB.UserRoleDao.CountUserIdsByRoleId(l.ctx, in.GetRoleId())
	}, func(offset int64, size int64) ([]int64, error) {
		return l.svcCtx.DB.UserRoleDao.PageUserIdsByRoleId(l.ctx, offset, size, in.GetRoleId())
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
				l.Logger.Errorf("查询角色关联的用户ID失败, roleId: %d, : %v", in.GetRoleId(), perr)
				return
			}

			source <- items
		}
	}, func(item []int64, writer mr.Writer[[]int64], cancel func(error)) {
		writer.Write(item)
	}, func(pipe <-chan []int64, cancel func(error)) {
		for items := range pipe {
			rerr := l.svcCtx.Cache.RoleCache.ClearRoles(l.ctx, role.SysType, items...)
			if rerr != nil {
				l.Logger.Errorf("清除用户角色缓存失败: %v", rerr)
				cancel(rerr)
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

	// 删除角色
	err = l.svcCtx.DB.RoleDao.Delete(l.ctx, in.GetRoleId())
	if err != nil {
		return nil, err
	}

	return &pb.RemoveRoleOutput{}, nil
}
