package gormx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/og-saas/framework/pkg/consts"
	"github.com/og-saas/framework/pkg/contextkey"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const dbOperation = "DB"

type GormLogger struct {
	logger.Config
}

func NewLogger(cfg logger.Config) *GormLogger {
	return &GormLogger{
		Config: cfg,
	}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		logx.WithContext(ctx).Infow(fmt.Sprintf(msg, data...), logx.Field("tenant", getTenantId(ctx)))
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		logx.WithContext(ctx).Sloww(fmt.Sprintf(msg, data...), logx.Field("tenant", getTenantId(ctx)))
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		logx.WithContext(ctx).Errorw(fmt.Sprintf(msg, data...), logx.Field("tenant", getTenantId(ctx)))
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []logx.LogField{
		logx.Field("sql", sql),
		logx.Field("rows", rows),
		logx.Field("duration", float64(elapsed.Nanoseconds())/1e6),
		logx.Field("tenant", getTenantId(ctx)),
	}

	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		fields = append(fields, logx.Field("err", err))
		logx.WithContext(ctx).Errorw(dbOperation, fields...)
	case elapsed > l.SlowThreshold && l.SlowThreshold > 0 && l.LogLevel >= logger.Warn:
		logx.WithContext(ctx).Sloww(dbOperation, fields...)
	case l.LogLevel == logger.Info:
		logx.WithContext(ctx).Infow(dbOperation, fields...)
	}
}

func getTenantId(ctx context.Context) any {
	tenantId := contextkey.GetContext[any](ctx, contextkey.TenantKey)
	return lo.Ternary(tenantId != nil, tenantId, consts.Default)
}
