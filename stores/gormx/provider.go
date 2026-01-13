package gormx

type TenantConfigProvider interface {
	Load() (map[string]Config, error)
}

type ConfigTenantProvider struct {
	configs map[string]Config
}

func NewConfigTenantProvider(config map[string]Config) *ConfigTenantProvider {
	return &ConfigTenantProvider{
		configs: config,
	}
}

func (p *ConfigTenantProvider) Load() (map[string]Config, error) {
	return p.configs, nil
}
