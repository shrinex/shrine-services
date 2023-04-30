package logic

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errAccountNotFound = status.Error(codes.NotFound, "用户不存在")
	errPasswdMismatch  = status.Error(codes.FailedPrecondition, "密码错误")
)
