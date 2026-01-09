package contextkey

// 空 struct 类型安全 context key 不占用内存，避免冲突
type (
	tenantKeyType struct{}
)

// 全局 key 实例
var (
	TenantKey = tenantKeyType{}
)
