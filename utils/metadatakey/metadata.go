package metadatakey

import "context"

func WithFromCtx(ctx context.Context, key, val any) context.Context {
	return context.WithValue(ctx, key, val)
}

func GetFromCtx(ctx context.Context, key any) any {
	return ctx.Value(key)
}

func GetFromCtxZero[T any](ctx context.Context, key any) (T, bool) {
	if val, ok := GetFromCtx(ctx, key).(T); ok {
		return val, true
	}
	var zero T
	return zero, false
}
