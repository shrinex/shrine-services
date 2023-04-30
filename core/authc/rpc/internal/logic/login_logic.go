package logic

import (
	"context"
	"core/authc/proto/model"
	"core/authc/rpc/internal/svc"
	"core/authc/rpc/pb"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Login 用户登录
func (l *LoginLogic) Login(in *pb.LoginInput) (*pb.LoginOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	acc, err := l.svcCtx.DB.AccountDao.FindOneBySysTypeUsername(l.ctx, in.GetSysType(), in.GetUsername())
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errAccountNotFound
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(in.GetPassword()))
	if err != nil {
		return nil, errPasswdMismatch
	}

	var out = new(pb.LoginOutput)
	_ = copier.Copy(out, acc)
	return out, nil
}
