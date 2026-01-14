package cachex

import (
	"context"
	"fmt"
	"github.com/og-saas/framework/utils/tenant"
)

type KeyType interface {
	CacheKey | TenantCacheKey | string
}

type (
	CacheKey       string
	TenantCacheKey string
)

func (key TenantCacheKey) String(ctx context.Context, args ...any) string {
	return fmt.Sprintf(string(key), append([]any{tenant.GetTenantId(ctx)}, args...)...)
}

func (key CacheKey) String(args ...any) string {
	return fmt.Sprintf(string(key), args...)
}

func KeyString[T KeyType](ctx context.Context, key T, args ...any) string {
	switch any(key).(type) {
	case CacheKey:
		return CacheKey(key).String(args...)
	case TenantCacheKey:
		return TenantCacheKey(key).String(ctx, args...)
	default:
		return string(key)
	}
}
