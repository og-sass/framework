package contextkey

import "context"

// SetContext 设置 context 值
func SetContext[T any](ctx context.Context, key any, val T) context.Context {
	return context.WithValue(ctx, key, val)
}

// GetContext 获取 context 值，如果不存在返回零值
func GetContext[T any](ctx context.Context, key any) T {
	var zero T
	if ctx == nil {
		return zero
	}
	if v := ctx.Value(key); v != nil {
		if val, ok := v.(T); ok {
			return val
		}
	}
	return zero
}
