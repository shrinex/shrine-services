package leaf

import (
	"github.com/sony/sonyflake"
	"time"
)

func init() {
	// 设置为东八区
	time.Local, _ = time.LoadLocation("Asia/Shanghai")
}

type (
	MachineID uint16

	Snowflake interface {
		NextID() (int64, error)
		MustNextID() int64
	}

	snowflake struct {
		adaptee *sonyflake.Sonyflake
	}
)

const (
	Authc MachineID = iota + 1
	Authz
	Merchant
	Dtm
	Product
)

var (
	_ Snowflake = (*snowflake)(nil)

	// startTime is a special time to me
	startTime = time.Date(2011, time.September, 1,
		0, 0, 0, 0, time.Local)
)

func NewSnowflake(mid MachineID) Snowflake {
	return &snowflake{
		adaptee: sonyflake.NewSonyflake(sonyflake.Settings{
			StartTime: startTime,
			MachineID: func() (uint16, error) {
				return uint16(mid), nil
			},
		}),
	}
}

func (s *snowflake) NextID() (int64, error) {
	val, err := s.adaptee.NextID()
	if err != nil {
		return 0, err
	}

	return int64(val), nil
}

func (s *snowflake) MustNextID() int64 {
	val, err := s.NextID()
	if err != nil {
		panic(err)
	}

	return val
}
