package gormx

import (
	"context"
	"sync"

	"github.com/og-saas/framework/pkg/consts"
	"github.com/og-saas/framework/pkg/contextkey"
	"gorm.io/gorm"
)

var Engine = NewDBEngine()

type DBEngine struct {
	Mysql      *DBManager
	Postgres   *DBManager
	Clickhouse *DBManager
}

func NewDBEngine() *DBEngine {
	return &DBEngine{}
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
	tenant := contextkey.GetContext[string](ctx, contextkey.TenantKey)
	if db, ok := s.getClientForTenant(tenant); ok {
		return db
	}

	if db, ok := s.getClientForTenant(consts.Default); ok {
		return db
	}
	panic("gorm: database not initialized")
}

func (s *DBManager) getClientForTenant(tenant string) (*gorm.DB, bool) {
	v, ok := s.pool.Load(tenant)
	if !ok || v == nil {
		return nil, false
	}
	db, ok := v.(*gorm.DB)
	return db, ok
}

func Must(configs ...Config) {
	must(consts.Default, configs...)
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
}

func must(tenant string, configs ...Config) {
	if tenant == "" {
		panic("gorm: empty tenant")
	}
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

		mgr.pool.Store(tenant, db)
	}
}
