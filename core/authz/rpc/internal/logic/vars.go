package logic

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errParentMenuNotExists = status.Error(codes.NotFound, "父菜单不存在")
)

const (
	rootMenuID int64 = 0
)
