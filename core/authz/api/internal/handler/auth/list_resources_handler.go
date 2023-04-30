package auth

import (
	"net/http"
	"shrine/std/globals"

	"core/authz/api/internal/logic/auth"
	"core/authz/api/internal/svc"
	"core/authz/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListResourcesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListResourcesReq
		if err := httpx.Parse(r, &req); err != nil {
			globals.ErrorCtx(r.Context(), w, err)
			return
		}

		l := auth.NewListResourcesLogic(r.Context(), svcCtx)
		resp, err := l.ListResources(&req)
		if err != nil {
			globals.ErrorCtx(r.Context(), w, err)
		} else {
			globals.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
