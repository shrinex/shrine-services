package realms

import (
	"context"
	"core/authc/rpc/pb"
	"core/authc/rpc/service"
	"github.com/shrinex/shield/authc"
	"shrine/std/authx"
)

type (
	sysAwareToken struct {
		sysType  int64
		username string
		password string
	}

	authcRpcRealm struct {
		authcRpc service.Service
	}
)

var _ authc.Token = (*sysAwareToken)(nil)

func NewToken(sysType int64, username string, password string) authc.Token {
	return &sysAwareToken{
		sysType:  sysType,
		username: username,
		password: password,
	}
}

func (t *sysAwareToken) Principal() string {
	return t.username
}

func (t *sysAwareToken) Credentials() string {
	return t.password
}

var _ authc.Realm = (*authcRpcRealm)(nil)

func NewAuthcRpcRealm(authcRpc service.Service) authc.Realm {
	return &authcRpcRealm{authcRpc: authcRpc}
}

func (r *authcRpcRealm) Supports(token authc.Token) bool {
	_, ok := token.(*sysAwareToken)
	return ok
}

func (r *authcRpcRealm) LoadUserDetails(ctx context.Context, token authc.Token) (authc.UserDetails, error) {
	tk := token.(*sysAwareToken)
	input := &pb.LoginInput{
		SysType:  tk.sysType,
		Username: tk.Principal(),
		Password: tk.Credentials(),
	}

	output, err := r.authcRpc.Login(ctx, input)
	if err != nil {
		return nil, err
	}

	return &authx.UserDetails{
		AccountId: output.AccountId,
		Username:  output.Username,
		UserId:    output.UserId,
		ShopId:    output.ShopId,
		SysType:   output.SysType,
		IsAdmin:   output.IsAdmin,
	}, nil
}
