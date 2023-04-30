package authx

import "github.com/shrinex/shield/authc"

type (
	UserDetails struct {
		AccountId int64  `json:"account_id"` // 唯一ID
		Username  string `json:"username"`   // 用户名
		UserId    int64  `json:"user_id"`    // 关联的用户ID
		ShopId    int64  `json:"shop_id"`    // 所属店铺
		SysType   int64  `json:"sys_type"`   // 系统类型: 1-平台端,2-商家端,3-普通用户
		IsAdmin   int64  `json:"is_admin"`   // 是否是管理员
	}
)

func (u *UserDetails) Principal() string {
	return u.Username
}

var _ authc.UserDetails = (*UserDetails)(nil)
