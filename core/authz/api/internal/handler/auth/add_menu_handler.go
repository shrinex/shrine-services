package auth

import (
	"net/http"
	"shrine/std/globals"

	"core/authz/api/internal/logic/auth"
	"core/authz/api/internal/svc"
	"core/authz/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AddMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddMenuReq
		if err := httpx.Parse(r, &req); err != nil {
			globals.ErrorCtx(r.Context(), w, err)
			return
		}

		l := auth.NewAddMenuLogic(r.Context(), svcCtx)
		resp, err := l.AddMenu(&req)
		if err != nil {
			globals.ErrorCtx(r.Context(), w, err)
		} else {
			globals.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
