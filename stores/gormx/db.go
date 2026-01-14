package gormx

import (
	"fmt"
	"github.com/og-saas/framework/stores/gormx/plugin"
	"strings"
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var dialectorFuncMap = map[string]func(string) gorm.Dialector{
	DriverMysql:      mysql.Open,
	DriverPostgres:   postgres.Open,
	DriverClickHouse: clickhouse.Open,
}

func level(level string) logger.LogLevel {
	switch strings.ToLower(level) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}

func (c Config) NewDB() *gorm.DB {
	f, ok := dialectorFuncMap[c.Driver]
	if !ok {
		panic("gorm: unsupported driver")
	}

	db, err := gorm.Open(f(c.DSN), &gorm.Config{
		SkipDefaultTransaction:    c.SkipDefaultTransaction,
		DefaultTransactionTimeout: time.Duration(c.DefaultTransactionTimeout) * time.Second,
		Logger: NewLogger(logger.Config{
			SlowThreshold: time.Duration(c.SlowThreshold) * time.Millisecond,
			LogLevel:      level(c.LogLevel),
		}),
		PrepareStmt:          c.PrepareStmt,
		PrepareStmtMaxSize:   c.PrepareStmtMaxSize,
		PrepareStmtTTL:       time.Duration(c.PrepareStmtTTL) * time.Second,
		DisableAutomaticPing: c.DisableAutomaticPing,
		TranslateError:       c.TranslateError,
	})
	if err != nil {
		panic(err)
	}

	if len(c.Sources) > 0 || len(c.Replicas) > 0 {
		var sources, replicas []gorm.Dialector
		for _, dsn := range c.Sources {
			sources = append(sources, f(dsn))
		}
		for _, dsn := range c.Replicas {
			replicas = append(replicas, f(dsn))
		}

		if err = db.Use(dbresolver.Register(dbresolver.Config{
			Sources:           sources,
			Replicas:          replicas,
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		})); err != nil {
			panic(fmt.Errorf("gorm dbresolver error: %w", err))
		}
	}

	// 4. 配置连接池
	if sqlDB, err := db.DB(); err == nil {
		if c.MaxIdleTime > 0 {
			sqlDB.SetConnMaxIdleTime(time.Duration(c.MaxIdleTime) * time.Second)
		}
		if c.MaxOpenConn > 0 {
			sqlDB.SetMaxOpenConns(c.MaxOpenConn)
		}
		if c.MaxIdleConn > 0 {
			sqlDB.SetMaxIdleConns(c.MaxIdleConn)
		}
		if c.MaxLifetime > 0 {
			sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
		}
	} else {
		panic(fmt.Errorf("gorm DB() error: %w", err))
	}

	plugin.NewTenantPlugin(c.TenantDBName).Register(db)
	return db
}
