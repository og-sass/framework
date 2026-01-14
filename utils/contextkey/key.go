package contextkey

type (
	tenantKeyType struct{}
	skipTenantKey struct{}
)

func (tenantKeyType) Name() string {
	return "tenant_id"
}

// 全局 key 实例
var (
	TenantKey     = tenantKeyType{}
	SkipTenantKey = skipTenantKey{}
)
