package redisx

import (
	"context"
	"fmt"
	"sync"

	"github.com/og-saas/framework/pkg/consts"
	"github.com/og-saas/framework/pkg/contextkey"
	"github.com/redis/go-redis/v9"
)

var Engine RDBEngine

type RDBEngine struct {
	pool sync.Map
}

func Must(c Config) {
	must(consts.Default, c)
}

func MustTenant(providers ...TenantConfigProvider) {
	for _, p := range providers {
		configMap, err := p.Load()
		if err != nil {
			panic(err)
		}

		for tenant, cfg := range configMap {
			must(tenant, cfg)
		}
	}
	fmt.Println("rdb: database initialized")
}

func must(tenant string, cfg Config) {
	if tenant == "" {
		panic("rdb: empty tenant")
	}

	rdb := cfg.newRdb()
	if rdb == nil {
		panic("rdb init failed")
	}
	Engine.pool.Store(tenant, rdb)
}

func (c Config) newRdb() (rdb redis.UniversalClient) {
	rdb = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:      c.Addrs,
		Username:   c.Username,
		Password:   c.Password,
		MasterName: c.MasterName,
		DB:         c.DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	if c.Debug {
		rdb.AddHook(DebugHook{})
	}

	return rdb
}

func (e *RDBEngine) RDB(ctx context.Context) redis.UniversalClient {

	tenant := contextkey.GetContext[string](ctx, contextkey.TenantKey)

	if rdb, ok := e.getClientForTenant(tenant); ok {
		return rdb
	}

	if rdb, ok := e.getClientForTenant(consts.Default); ok {
		return rdb
	}

	return nil

}

// getClientForTenant 从 pool 安全获取 redis client
func (e *RDBEngine) getClientForTenant(tenant string) (redis.UniversalClient, bool) {
	v, ok := e.pool.Load(tenant)
	if !ok || v == nil {
		return nil, false
	}
	rdb, ok := v.(redis.UniversalClient)
	return rdb, ok
}
