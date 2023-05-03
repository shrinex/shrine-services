package shop

import (
	"net/http"
	"shrine/std/globals"

	"biz/platform/api/internal/logic/shop"
	"biz/platform/api/internal/svc"
	"biz/platform/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateShopReq
		if err := httpx.Parse(r, &req); err != nil {

			globals.ErrorCtx(r.Context(), w, err)
			return
		}

		l := shop.NewCreateShopLogic(r.Context(), svcCtx)
		resp, err := l.CreateShop(&req)
		if err != nil {
			globals.ErrorCtx(r.Context(), w, err)
		} else {
			globals.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
