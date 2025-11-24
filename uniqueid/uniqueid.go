package uniqueid

import (
	"time"

	"github.com/sony/sonyflake"
)

var (
	flake *sonyflake.Sonyflake
)

func init() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: time.Date(1997, 1, 14, 0, 0, 0, 0, time.UTC),
	})
	if flake == nil {
		panic("sony flake init failed")
	}
}

// GenId 生成一个唯一的雪花ID
func GenId() (id uint64, err error) {
	id, err = flake.NextID()
	return
}
