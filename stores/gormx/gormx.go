package gormx

import (
	"context"
	"github.com/og-saas/framework/utils/tenant"
	"github.com/spf13/cast"
	"sync"

	"gorm.io/gorm"
)

var (
	Engine *DBEngine
	once   sync.Once
)

type DBEngine struct {
	Mysql      *DBManager
	Postgres   *DBManager
	Clickhouse *DBManager
}

type DBManager struct {
	pool sync.Map
}

func (e *DBEngine) DB(drivers ...string) *DBManager {
	var driver string
	if len(drivers) == 0 || drivers[0] == "" {
		driver = DriverMysql // 默认 MySQL
	} else {
		driver = drivers[0]
	}
	switch driver {
	case DriverMysql:
		return e.Mysql
	case DriverPostgres:
		return e.Postgres
	case DriverClickHouse:
		return e.Clickhouse
	}
	panic("gorm: unknown driver")
}

func (s *DBManager) WithContext(ctx context.Context) *gorm.DB {
	if db, ok := s.getClientForTenant(tenant.GetTenantId(ctx)); ok {
		return db.WithContext(ctx)
	}

	if db, ok := s.getClientForTenant(tenant.Default); ok {
		return db.WithContext(ctx)
	}
	panic("gorm: database not initialized")
}

func (s *DBManager) getClientForTenant(tenantId int64) (*gorm.DB, bool) {
	v, ok := s.pool.Load(tenantId)
	if !ok || v == nil {
		return nil, false
	}
	db, ok := v.(*gorm.DB)
	return db, ok
}

func Must(configs ...Config) {
	must(tenant.Default, configs...)
}

func MustTenant(providers ...TenantConfigProvider) {
	for _, p := range providers {
		configMap, err := p.Load()
		if err != nil {
			panic(err)
		}

		for key, val := range configMap {
			tenantId := cast.ToInt64(key)
			must(tenantId, val)
		}
	}
}

func must(tenantId int64, configs ...Config) {
	once.Do(func() {
		Engine = &DBEngine{}
	})
	if len(configs) == 0 {
		panic("gorm: empty config")
	}

	for _, cfg := range configs {
		db := cfg.NewDB()
		if db == nil {
			panic("gorm: db init failed")
		}

		var mgr *DBManager

		switch cfg.Driver {
		case DriverMysql:
			if Engine.Mysql == nil {
				Engine.Mysql = &DBManager{}
			}
			mgr = Engine.Mysql

		case DriverPostgres:
			if Engine.Postgres == nil {
				Engine.Postgres = &DBManager{}
			}
			mgr = Engine.Postgres

		case DriverClickHouse:
			if Engine.Clickhouse == nil {
				Engine.Clickhouse = &DBManager{}
			}
			mgr = Engine.Clickhouse

		default:
			panic("gorm: unknown driver")
		}

		mgr.pool.Store(tenantId, db)
	}
}
