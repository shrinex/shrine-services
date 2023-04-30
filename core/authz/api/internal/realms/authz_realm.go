package realms

import (
	"context"
	"core/authz/rpc/pb"
	"core/authz/rpc/service"
	"github.com/shrinex/shield-web/pattern"
	"github.com/shrinex/shield/authc"
	"github.com/shrinex/shield/authz"
	"shrine/std/authx"
	"shrine/std/utils/slices"
)

type (
	authzRpcRealm struct {
		authzRpc service.Service
	}

	authzRole struct {
		roleId int64
		name   string
	}

	authzAuthority struct {
		name    string
		method  string
		pattern string
	}
)

var (
	_ authz.Role      = (*authzRole)(nil)
	_ authz.Authority = (*authzAuthority)(nil)

	antMatcher = pattern.NewMatcher()
)

func NewAuthzRole(roleId int64, name string) authz.Role {
	return &authzRole{roleId: roleId, name: name}
}

func (r *authzRole) Desc() string {
	return r.name
}

func (r *authzRole) Implies(role authz.Role) bool {
	that, ok := role.(*authzRole)
	if !ok {
		return false
	}

	return r.roleId == that.roleId
}

func NewAuthzAuthority(method string, pattern string) authz.Authority {
	return &authzAuthority{name: "request", method: method, pattern: pattern}
}

func NewAuthzAuthorityWithName(name string, method string, pattern string) authz.Authority {
	return &authzAuthority{name: name, method: method, pattern: pattern}
}

func (a *authzAuthority) Desc() string {
	return a.name
}

func (a *authzAuthority) Implies(authority authz.Authority) bool {
	that, ok := authority.(*authzAuthority)
	if !ok {
		return false
	}

	return a.method == that.method && antMatcher.Matches(a.pattern, that.pattern)
}

func NewAuthzRpcRealm(authzRpc service.Service) authz.Realm {
	return &authzRpcRealm{authzRpc: authzRpc}
}

func (r *authzRpcRealm) LoadRoles(ctx context.Context, userDetails authc.UserDetails) ([]authz.Role, error) {
	user := userDetails.(*authx.UserDetails)
	input := &pb.ListRolesInput{
		UserId:  user.UserId,
		SysType: user.SysType,
		IsAdmin: user.IsAdmin,
	}

	output, err := r.authzRpc.ListRoles(ctx, input)
	if err != nil {
		return slices.Empty[authz.Role](), err
	}

	var roles []authz.Role
	for _, r := range output.GetRoles() {
		roles = append(roles, NewAuthzRole(r.GetRoleId(), r.GetName()))
	}

	return roles, nil
}

func (r *authzRpcRealm) LoadAuthorities(ctx context.Context, userDetails authc.UserDetails) ([]authz.Authority, error) {
	user := userDetails.(*authx.UserDetails)
	input := &pb.ListResourcesInput{
		UserId:  user.UserId,
		SysType: user.SysType,
		IsAdmin: user.IsAdmin,
	}

	output, err := r.authzRpc.ListResources(ctx, input)
	if err != nil {
		return slices.Empty[authz.Authority](), err
	}

	var authorities []authz.Authority
	for _, r := range output.GetResources() {
		authorities = append(authorities, NewAuthzAuthorityWithName(r.GetName(), r.GetMethod(), r.GetPattern()))
	}

	return authorities, nil
}
