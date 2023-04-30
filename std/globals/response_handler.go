package globals

import (
	"context"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type response struct {
	Code    int32  `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
	Data    any    `json:"data"`    // 响应体
}

func ErrorCtx(ctx context.Context, w http.ResponseWriter, err error) {
	if s, ok := status.FromError(err); ok {
		httpx.OkJsonCtx(ctx, w, &response{
			Code:    int32(s.Code()),
			Message: s.Message(),
		})
		return
	}

	httpx.OkJsonCtx(ctx, w, &response{
		Code:    int32(codes.Internal),
		Message: err.Error(),
	})
}

func OkJsonCtx(ctx context.Context, w http.ResponseWriter, v any) {
	httpx.OkJsonCtx(ctx, w, &response{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    v,
	})
}
