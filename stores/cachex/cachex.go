package cachex

import (
	"sync"

	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"
	"saas.og.game/framework/stores/cachex/config"
)

var Engine *rockscache.Client
var once sync.Once

func Must(c config.Config, rdb redis.UniversalClient) {
	once.Do(func() {
		if Engine == nil {
			Engine = NewEngine(c, rdb)
		}
	})
}

// NewEngine 创建一个缓存引擎
func NewEngine(c config.Config, rdb redis.UniversalClient) *rockscache.Client {
	options := rockscache.NewDefaultOptions()
	options.StrongConsistency = c.StrongConsistency
	options.DisableCacheRead = c.DisableCacheRead
	return rockscache.NewClient(rdb, options)
}
