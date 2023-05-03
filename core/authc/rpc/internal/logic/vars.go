package logic

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errUserExistsDesc    = "用户已存在"
	errAccountExistsDesc = "账号已存在"

	errUserExists      = status.Error(codes.AlreadyExists, errUserExistsDesc)
	errAccountExists   = status.Error(codes.AlreadyExists, errAccountExistsDesc)
	errAccountNotFound = status.Error(codes.NotFound, "账号不存在")
	errPasswdMismatch  = status.Error(codes.FailedPrecondition, "密码错误")
)
