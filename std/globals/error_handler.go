package globals

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type errorEnvelope struct {
	Code    int32  `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
}

func ErrorHandler(_ context.Context, err error) (int, any) {
	if s, ok := status.FromError(err); ok {
		return http.StatusOK, &errorEnvelope{
			Code:    int32(s.Code()),
			Message: s.Message(),
		}
	}

	return http.StatusOK, &errorEnvelope{
		Code:    int32(codes.Internal),
		Message: err.Error(),
	}
}
