package globals

import (
	"context"
	"github.com/pkg/errors"
	"github.com/shrinex/shield/authc"
	"github.com/shrinex/shield/semgt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type Response struct {
	Code    int32  `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
	Data    any    `json:"data"`    // 响应体
}

func ErrorCtx(ctx context.Context, w http.ResponseWriter, err error) {
	if s, ok := status.FromError(err); ok {
		httpx.OkJsonCtx(ctx, w, &Response{
			Code:    int32(s.Code()),
			Message: s.Message(),
		})
		return
	}

	httpx.OkJsonCtx(ctx, w, evalResponse(err))
}

func OkJsonCtx(ctx context.Context, w http.ResponseWriter, v any) {
	httpx.OkJsonCtx(ctx, w, &Response{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    v,
	})
}

func evalResponse(err error) *Response {
	if errors.Is(err, authc.ErrInvalidToken) {
		return &Response{
			Code:    http.StatusBadRequest,
			Message: "token格式不正确",
		}
	}

	if errors.Is(err, authc.ErrUnauthenticated) {
		return &Response{
			Code:    http.StatusUnauthorized,
			Message: "请先登录",
		}
	}

	if errors.Is(err, semgt.ErrExpired) {
		return &Response{
			Code:    int32(SessionExpired),
			Message: "会话已过期，请重新登录",
		}
	}

	if errors.Is(err, semgt.ErrReplaced) {
		return &Response{
			Code:    int32(SessionReplaced),
			Message: "当前账号已在其它设备登录",
		}
	}

	if errors.Is(err, semgt.ErrOverflow) {
		return &Response{
			Code:    int32(SessionOverflow),
			Message: "会话已超限，请重新登录",
		}
	}

	return &Response{
		Code:    int32(codes.Internal),
		Message: err.Error(),
	}
}
