package logic

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errUserNotExists = status.Error(codes.NotFound, "用户不存在")
)
