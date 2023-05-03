package logic

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errShopExistsDesc = "店铺已存在"

	errShopExists = status.Error(codes.AlreadyExists, errShopExistsDesc)
)
