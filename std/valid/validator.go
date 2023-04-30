package valid

import (
	"github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func init() {
	httpx.SetValidator(&tagExprBased{
		Validator: validator.New("valid"),
	})
}

type (
	tagExprBased struct {
		*validator.Validator
	}
)

var _ httpx.Validator = (*tagExprBased)(nil)

func (v *tagExprBased) Validate(_ *http.Request, data any) error {
	return v.Validator.Validate(data)
}
