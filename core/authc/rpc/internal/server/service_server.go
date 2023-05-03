// Code generated by goctl. DO NOT EDIT.
// Source: main.proto

package server

import (
	"context"

	"core/authc/rpc/internal/logic"
	"core/authc/rpc/internal/svc"
	"core/authc/rpc/pb"
)

type ServiceServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedServiceServer
}

func NewServiceServer(svcCtx *svc.ServiceContext) *ServiceServer {
	return &ServiceServer{
		svcCtx: svcCtx,
	}
}

// Login 用户登录
func (s *ServiceServer) Login(ctx context.Context, in *pb.LoginInput) (*pb.LoginOutput, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(in)
}

// Register 用户注册
func (s *ServiceServer) Register(ctx context.Context, in *pb.RegisterInput) (*pb.RegisterOutput, error) {
	l := logic.NewRegisterLogic(ctx, s.svcCtx)
	return l.Register(in)
}

// RegisterConfirm 用户注册确认
func (s *ServiceServer) RegisterConfirm(ctx context.Context, in *pb.RegisterInput) (*pb.RegisterOutput, error) {
	l := logic.NewRegisterConfirmLogic(ctx, s.svcCtx)
	return l.RegisterConfirm(in)
}

// RegisterCancel 用户注册回滚
func (s *ServiceServer) RegisterCancel(ctx context.Context, in *pb.RegisterInput) (*pb.RegisterOutput, error) {
	l := logic.NewRegisterCancelLogic(ctx, s.svcCtx)
	return l.RegisterCancel(in)
}
