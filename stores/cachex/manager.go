package cachex

import (
	"context"
	"time"
)

// Fetch2
// Key 只支持 CacheKey | TenantCacheKey | string 类型
// key为string 不会带入参数
// 如果是 TenantCacheKey, 不用传递tenant参数
func Fetch2[T KeyType](ctx context.Context, key T, expire time.Duration, fn func() (string, error), args ...any) (string, error) {
	return Engine(ctx).Fetch2(ctx, KeyString(ctx, key, args...), expire, fn)
}

// TagAsDeleted2 Key 只支持 CacheKey | TenantCacheKey | string 类型, string不会带入参数
func TagAsDeleted2[T KeyType](ctx context.Context, key T, args ...any) error {
	return Engine(ctx).TagAsDeleted2(ctx, KeyString(ctx, key, args...))
}
