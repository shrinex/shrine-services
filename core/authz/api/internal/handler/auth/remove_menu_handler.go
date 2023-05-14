package auth

import (
	"net/http"
	"shrine/std/globals"

	"core/authz/api/internal/logic/auth"
	"core/authz/api/internal/svc"
	"core/authz/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RemoveMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RemoveMenuReq
		if err := httpx.Parse(r, &req); err != nil {
			globals.ErrorCtx(r.Context(), w, err)
			return
		}

		l := auth.NewRemoveMenuLogic(r.Context(), svcCtx)
		resp, err := l.RemoveMenu(&req)
		if err != nil {
			globals.ErrorCtx(r.Context(), w, err)
		} else {
			globals.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
