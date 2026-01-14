package plugin

import (
	"github.com/og-saas/framework/utils/tenant"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// ==================== 插件 ====================

type TenantPlugin struct {
	TenantDBName string
}

// NewTenantPlugin 创建插件实例
func NewTenantPlugin(tenantDbName string) *TenantPlugin {
	return &TenantPlugin{
		TenantDBName: tenantDbName,
	}
}

// Register 注册到 GORM DB
func (p *TenantPlugin) Register(db *gorm.DB) {
	_ = db.Callback().Query().Before("gorm:query").Register("tenant:query", p.queryHook)
	_ = db.Callback().Create().Before("gorm:create").Register("tenant:create", p.createHook)
	_ = db.Callback().Update().Before("gorm:update").Register("tenant:update", p.updateHook)
	_ = db.Callback().Delete().Before("gorm:delete").Register("tenant:delete", p.deleteHook)
}

// ==================== Hook 核心 ====================

// Create Hook
func (p *TenantPlugin) createHook(db *gorm.DB) {
	var (
		skip  bool
		field *schema.Field
	)

	if field, skip = p.before(db); skip {
		return
	}

	_ = field.Set(db.Statement.Context, db.Statement.ReflectValue, tenant.GetTenantId(db.Statement.Context))
}

// queryHook
func (p *TenantPlugin) queryHook(db *gorm.DB) {
	p.tenantCond(db)
}

// updateHook
func (p *TenantPlugin) updateHook(db *gorm.DB) {
	p.tenantCond(db)
}

// deleteHook
func (p *TenantPlugin) deleteHook(db *gorm.DB) {
	p.tenantCond(db)
}

func (p *TenantPlugin) before(db *gorm.DB) (*schema.Field, bool) {
	if db.Statement.Schema == nil {
		return nil, true
	}

	// 跳过
	if tenant.IsSkipTenant(db.Statement.Context) {
		return nil, true
	}

	field, ok := db.Statement.Schema.FieldsByDBName[p.TenantDBName]
	return field, !ok
}

func (p *TenantPlugin) tenantCond(db *gorm.DB) {
	var (
		skip  bool
		field *schema.Field
	)

	if field, skip = p.before(db); skip {
		return
	}

	db.Statement.AddClause(clause.Where{
		Exprs: []clause.Expression{
			clause.Eq{Column: field.DBName, Value: tenant.GetTenantId(db.Statement.Context)},
		},
	})
}
