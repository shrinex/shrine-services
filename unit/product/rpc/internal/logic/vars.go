package logic

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	rootCategoryId   = 0
	maxCategoryLevel = 3
)

var (
	errCategoryNotFound       = status.Error(codes.NotFound, "分类不存在")
	errParentCategoryNotFound = status.Error(codes.NotFound, "父级分类不存在")
	errCategoryLevelOverflow  = status.Error(codes.InvalidArgument, "最多支持三级分类")
)
