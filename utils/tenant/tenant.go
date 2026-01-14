package tenant

import (
	"context"
	"github.com/og-saas/framework/utils/contextkey"
	"github.com/spf13/cast"
)

const Default int64 = 0

// GetTenantId 获取租户ID 如果跳过租户则返回0
func GetTenantId(ctx context.Context) int64 {
	if IsSkipTenant(ctx) {
		return Default
	}
	val := contextkey.GetContext[any](ctx, contextkey.TenantKey)
	return cast.ToInt64(val)
}

// SetTenantId 设置租户
func SetTenantId(ctx context.Context, tenantId int64) context.Context {
	return contextkey.SetContext(ctx, contextkey.TenantKey, tenantId)
}

// SkipTenant 上下文设置跳过租户
func SkipTenant(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextkey.SkipTenantKey, true)
}

// IsSkipTenant 检查是否跳过租户
func IsSkipTenant(ctx context.Context) bool {
	return contextkey.GetContext[bool](ctx, contextkey.SkipTenantKey)
}
