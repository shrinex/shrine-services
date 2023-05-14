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

type RemoveMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveMenuLogic {
	return &RemoveMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RemoveMenu 删除菜单
func (l *RemoveMenuLogic) RemoveMenu(in *pb.RemoveMenuInput) (*pb.RemoveMenuOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	menu, err := l.svcCtx.DB.MenuDao.FindOne(l.ctx, in.GetMenuId())
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return &pb.RemoveMenuOutput{}, nil
		}
		return nil, err
	}

	// 删除 角色->菜单 缓存
	var perr error
	pager := iter.NewPager(func() (int64, error) {
		return l.svcCtx.DB.RoleMenuDao.CountRoleIdsByMenuId(l.ctx, in.GetMenuId())
	}, func(offset int64, size int64) ([]int64, error) {
		return l.svcCtx.DB.RoleMenuDao.PageRoleIdsByMenuId(l.ctx, offset, size, in.GetMenuId())
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
				l.Logger.Errorf("查询菜单关联的角色ID失败, menuId: %d, err: %v", in.GetMenuId(), perr)
				return
			}

			source <- items
		}
	}, func(items []int64, writer mr.Writer[[]int64], cancel func(error)) {
		writer.Write(items)
	}, func(pipe <-chan []int64, cancel func(error)) {
		for items := range pipe {
			rerr := l.svcCtx.Cache.MenuCache.ClearMenus(l.ctx, menu.SysType, items...)
			if rerr != nil {
				l.Logger.Errorf("清除用户菜单缓存失败: %v", rerr)
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

	// 删除菜单
	err = l.svcCtx.DB.MenuDao.Delete(l.ctx, in.GetMenuId())
	if err != nil {
		return nil, err
	}

	return &pb.RemoveMenuOutput{}, nil
}
