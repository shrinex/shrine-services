package globals

import "github.com/pkg/errors"

var (
	ErrForbidden = errors.New("权限不足")
)
