package redisx

import (
	"context"
	"github.com/og-saas/framework/utils/tenant"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	Engine *RDBEngine
	once   sync.Once
)

type RDBEngine struct {
	pool sync.Map
}

func Must(c Config) {
	must(tenant.Default, c)
}

func MustTenant(providers ...TenantConfigProvider) {
	for _, p := range providers {
		configMap, err := p.Load()
		if err != nil {
			panic(err)
		}

		for key, val := range configMap {
			must(key, val)
		}
	}
}

func must(tenant int64, cfg Config) {
	once.Do(func() {
		Engine = &RDBEngine{}
	})
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

	if rdb, ok := e.getClientForTenant(tenant.GetTenantId(ctx)); ok {
		return rdb
	}

	if rdb, ok := e.getClientForTenant(tenant.Default); ok {
		return rdb
	}

	return nil

}

// getClientForTenant 从 pool 安全获取 redis client
func (e *RDBEngine) getClientForTenant(tenantId int64) (redis.UniversalClient, bool) {
	v, ok := e.pool.Load(tenantId)
	if !ok || v == nil {
		return nil, false
	}
	rdb, ok := v.(redis.UniversalClient)
	return rdb, ok
}

func (e *RDBEngine) Map() map[int64]redis.UniversalClient {
	m := make(map[int64]redis.UniversalClient)
	e.pool.Range(func(key, value interface{}) bool {
		m[key.(int64)] = value.(redis.UniversalClient)
		return true
	})
	return m
}
