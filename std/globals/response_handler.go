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

	httpx.OkJsonCtx(ctx, w, &Response{
		Code:    int32(codes.Internal),
		Message: evalMessage(err),
	})
}

func OkJsonCtx(ctx context.Context, w http.ResponseWriter, v any) {
	httpx.OkJsonCtx(ctx, w, &Response{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    v,
	})
}

func evalMessage(err error) string {
	if errors.Is(err, authc.ErrInvalidToken) {
		return "token格式不正确"
	}

	if errors.Is(err, authc.ErrUnauthenticated) {
		return "请先登录"
	}

	if errors.Is(err, semgt.ErrExpired) {
		return "会话已过期，请重新登录"
	}

	if errors.Is(err, semgt.ErrReplaced) {
		return "当前账号已在其它设备登录"
	}

	if errors.Is(err, semgt.ErrOverflow) {
		return "会话已超限，请重新登录"
	}

	return err.Error()
}
