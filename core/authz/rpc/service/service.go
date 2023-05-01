// Code generated by goctl. DO NOT EDIT.
// Source: main.proto

package service

import (
	"context"

	"core/authz/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AddMenuInput              = pb.AddMenuInput
	AddMenuOutput             = pb.AddMenuOutput
	AddResourceGroupInput     = pb.AddResourceGroupInput
	AddResourceGroupOutput    = pb.AddResourceGroupOutput
	AddResourceInput          = pb.AddResourceInput
	AddResourceOutput         = pb.AddResourceOutput
	AddRoleInput              = pb.AddRoleInput
	AddRoleOutput             = pb.AddRoleOutput
	ListResourcesInput        = pb.ListResourcesInput
	ListResourcesOutput       = pb.ListResourcesOutput
	ListRolesInput            = pb.ListRolesInput
	ListRolesOutput           = pb.ListRolesOutput
	Menu                      = pb.Menu
	PageMenusInput            = pb.PageMenusInput
	PageMenusOutput           = pb.PageMenusOutput
	PageResourcesGroupsInput  = pb.PageResourcesGroupsInput
	PageResourcesGroupsOutput = pb.PageResourcesGroupsOutput
	RemoveMenuInput           = pb.RemoveMenuInput
	RemoveMenuOutput          = pb.RemoveMenuOutput
	RemoveResourceGroupInput  = pb.RemoveResourceGroupInput
	RemoveResourceGroupOutput = pb.RemoveResourceGroupOutput
	RemoveResourceInput       = pb.RemoveResourceInput
	RemoveResourceOutput      = pb.RemoveResourceOutput
	RemoveRoleInput           = pb.RemoveRoleInput
	RemoveRoleOutput          = pb.RemoveRoleOutput
	Resource                  = pb.Resource
	ResourceGroup             = pb.ResourceGroup
	Role                      = pb.Role

	Service interface {
		// AddRole 添加角色
		AddRole(ctx context.Context, in *AddRoleInput, opts ...grpc.CallOption) (*AddRoleOutput, error)
		// AddResource 添加资源
		AddResource(ctx context.Context, in *AddResourceInput, opts ...grpc.CallOption) (*AddResourceOutput, error)
		// AddResourceGroup 添加资源分组
		AddResourceGroup(ctx context.Context, in *AddResourceGroupInput, opts ...grpc.CallOption) (*AddResourceGroupOutput, error)
		// AddMenu 添加菜单
		AddMenu(ctx context.Context, in *AddMenuInput, opts ...grpc.CallOption) (*AddMenuOutput, error)
		// RemoveRole 删除角色
		RemoveRole(ctx context.Context, in *RemoveRoleInput, opts ...grpc.CallOption) (*RemoveRoleOutput, error)
		// RemoveResource 删除资源
		RemoveResource(ctx context.Context, in *RemoveResourceInput, opts ...grpc.CallOption) (*RemoveResourceOutput, error)
		// RemoveResourceGroup 删除资源分组
		RemoveResourceGroup(ctx context.Context, in *RemoveResourceGroupInput, opts ...grpc.CallOption) (*RemoveResourceGroupOutput, error)
		// RemoveMenu 删除菜单
		RemoveMenu(ctx context.Context, in *RemoveMenuInput, opts ...grpc.CallOption) (*RemoveMenuOutput, error)
		// ListRoles 查询用户拥有的角色列表
		ListRoles(ctx context.Context, in *ListRolesInput, opts ...grpc.CallOption) (*ListRolesOutput, error)
		// ListResources 查询用户拥有的资源列表
		ListResources(ctx context.Context, in *ListResourcesInput, opts ...grpc.CallOption) (*ListResourcesOutput, error)
		// PageResourcesGroups 分页查询系统中的资源分组
		PageResourcesGroups(ctx context.Context, in *PageResourcesGroupsInput, opts ...grpc.CallOption) (*PageResourcesGroupsOutput, error)
		// PageMenus 分页查询系统中的菜单列表
		PageMenus(ctx context.Context, in *PageMenusInput, opts ...grpc.CallOption) (*PageMenusOutput, error)
	}

	defaultService struct {
		cli zrpc.Client
	}
)

func NewService(cli zrpc.Client) Service {
	return &defaultService{
		cli: cli,
	}
}

// AddRole 添加角色
func (m *defaultService) AddRole(ctx context.Context, in *AddRoleInput, opts ...grpc.CallOption) (*AddRoleOutput, error) {
	client := pb.NewServiceClient(m.cli.Conn())
	return client.AddRole(ctx, in, opts...)
}

// AddResource 添加资源
func (m *defaultService) AddResource(ctx context.Context, in *AddResourceInput, opts ...grpc.CallOption) (*AddResourceOutput, error) {
	client := pb.NewServiceClient(m.cli.Conn())
	return client.AddResource(ctx, in, opts...)
}

// AddResourceGroup 添加资源分组
func (m *defaultService) AddResourceGroup(ctx context.Context, in *AddResourceGroupInput, opts ...grpc.CallOption) (*AddResourceGroupOutput, error) {
	client := pb.NewServiceClient(m.cli.Conn())
	return client.AddResourceGroup(ctx, in, opts...)
}

// AddMenu 添加菜单
func (m *defaultService) AddMenu(ctx context.Context, in *AddMenuInput, opts ...grpc.CallOption) (*AddMenuOutput, error) {
	client := pb.NewServiceClient(m.cli.Conn())
	return client.AddMenu(ctx, in, opts...)
}

// RemoveRole 删除角色
func (m *defaultService) RemoveRole(ctx context.Context, in *RemoveRoleInput, opts ...grpc.CallOption) (*RemoveRoleOutput, error) {
	client := pb.NewServiceClient(m.cli.Conn())
	return client.RemoveRole(ctx, in, opts...)
}

// RemoveResource 删除资源
func (m *defaultService) RemoveResource(ctx context.Context, in *RemoveResourceInput, opts ...grpc.CallOption) (*RemoveResourceOutput, error) {
	client := pb.NewServiceClient(m.cli.Conn())
	return client.RemoveResource(ctx, in, opts...)
}

// RemoveResourceGroup 删除资源分组
func (m *defaultService) RemoveResourceGroup(ctx context.Context, in *RemoveResourceGroupInput, opts ...grpc.CallOption) (*RemoveResourceGroupOutput, error) {
	client := pb.NewServiceClient(m.cli.Conn())
	return client.RemoveResourceGroup(ctx, in, opts...)
}

// RemoveMenu 删除菜单
func (m *defaultService) RemoveMenu(ctx context.Context, in *RemoveMenuInput, opts ...grpc.CallOption) (*RemoveMenuOutput, error) {
	client := pb.NewServiceClient(m.cli.Conn())
	return client.RemoveMenu(ctx, in, opts...)
}

// ListRoles 查询用户拥有的角色列表
func (m *defaultService) ListRoles(ctx context.Context, in *ListRolesInput, opts ...grpc.CallOption) (*ListRolesOutput, error) {
	client := pb.NewServiceClient(m.cli.Conn())
	return client.ListRoles(ctx, in, opts...)
}

// ListResources 查询用户拥有的资源列表
func (m *defaultService) ListResources(ctx context.Context, in *ListResourcesInput, opts ...grpc.CallOption) (*ListResourcesOutput, error) {
	client := pb.NewServiceClient(m.cli.Conn())
	return client.ListResources(ctx, in, opts...)
}

// PageResourcesGroups 分页查询系统中的资源分组
func (m *defaultService) PageResourcesGroups(ctx context.Context, in *PageResourcesGroupsInput, opts ...grpc.CallOption) (*PageResourcesGroupsOutput, error) {
	client := pb.NewServiceClient(m.cli.Conn())
	return client.PageResourcesGroups(ctx, in, opts...)
}

// PageMenus 分页查询系统中的菜单列表
func (m *defaultService) PageMenus(ctx context.Context, in *PageMenusInput, opts ...grpc.CallOption) (*PageMenusOutput, error) {
	client := pb.NewServiceClient(m.cli.Conn())
	return client.PageMenus(ctx, in, opts...)
}