package conv

import (
	"github.com/jinzhu/copier"
	"time"
)

func NewTime2int64() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: time.Time{},
		DstType: int64(0),
		Fn: func(src any) (any, error) {
			return src.(time.Time).Unix(), nil
		},
	}
}
