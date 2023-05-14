package auth

import (
	"context"
	"core/authz/api/internal/svc"
	"core/authz/api/internal/types"
	"core/authz/rpc/pb"
	"github.com/jinzhu/copier"
	"shrine/std/authx"
	"shrine/std/utils/slices"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenusLogic {
	return &ListMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMenusLogic) ListMenus(_ *types.ListMenusReq) (resp *types.ListMenusResp, err error) {
	menus := slices.Empty[*types.MenuNode]()
	userDetails, err := l.svcCtx.Subject.UserDetails(l.ctx)
	if err != nil {
		resp = &types.ListMenusResp{
			Menus: menus,
		}
		return
	}

	user := userDetails.(*authx.UserDetails)
	input := &pb.ListMenusInput{
		UserId:  user.UserId,
		SysType: user.SysType,
		IsAdmin: user.IsAdmin,
	}

	output, err := l.svcCtx.AuthzRpc.ListMenus(l.ctx, input)
	if err != nil {
		resp = &types.ListMenusResp{
			Menus: menus,
		}
		return
	}

	resp = &types.ListMenusResp{
		Menus: l.treeify(output.GetMenus()),
	}
	return
}

func (l *ListMenusLogic) treeify(menus []*pb.Menu) []*types.MenuNode {
	nodes := slices.Map(menus, func(e *pb.Menu) (ret *types.MenuNode) {
		ret = new(types.MenuNode)
		_ = copier.Copy(ret, e)
		return
	})

	lookup := slices.GroupingBy(nodes, func(e *types.MenuNode) int64 {
		return e.ParentId
	})

	for _, e := range nodes {
		if children, ok := lookup[e.MenuId]; ok {
			e.Children = children
		} else {
			e.Children = []*types.MenuNode{}
		}
	}

	return slices.Filter(nodes, func(e *types.MenuNode) bool {
		return e.ParentId == rootMenuId
	})
}
