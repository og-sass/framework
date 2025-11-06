package uniqueid

import (
	"time"

	"github.com/sony/sonyflake"
	"github.com/zeromicro/go-zero/core/logx"
)

var flake *sonyflake.Sonyflake

// startTime 1997-01-14 00:00:00
var startTime = time.Date(1997, 1, 14, 0, 0, 0, 0, time.UTC)

func init() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: startTime,
	})
}

func GenId() int64 {

	id, err := flake.NextID()
	if err != nil {
		logx.Severef("flake NextID failed with %s \n", err)
		panic(err)
	}

	return int64(id)
}
