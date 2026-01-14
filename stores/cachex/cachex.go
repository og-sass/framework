package cachex

import (
	"context"
	"github.com/dtm-labs/rockscache"
	"github.com/og-saas/framework/stores/redisx"
	"github.com/og-saas/framework/utils/tenant"
	"github.com/redis/go-redis/v9"
	"sync"
)

var engine sync.Map

func MustTenant(c Config, rdbEngine *redisx.RDBEngine) {
	for tenantId, client := range rdbEngine.Map() {
		engine.Store(tenantId, New(c, client))
	}
}

func Engine(ctx context.Context) *rockscache.Client {
	if rdb, ok := getClientForTenant(tenant.GetTenantId(ctx)); ok {
		return rdb
	}

	if rdb, ok := getClientForTenant(tenant.Default); ok {
		return rdb
	}

	return nil
}

// New 创建一个缓存
func New(c Config, rdb redis.UniversalClient) *rockscache.Client {
	options := rockscache.NewDefaultOptions()
	options.StrongConsistency = c.StrongConsistency
	options.DisableCacheRead = c.DisableCacheRead
	return rockscache.NewClient(rdb, options)
}

// getClientForTenant 从 pool 安全获取 client
func getClientForTenant(tenantId int64) (*rockscache.Client, bool) {
	v, ok := engine.Load(tenantId)
	if !ok || v == nil {
		return nil, false
	}
	cli, ok := v.(*rockscache.Client)
	return cli, ok
}
