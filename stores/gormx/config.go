package gormx

const (
	DriverMysql      = "mysql"
	DriverPostgres   = "postgres"
	DriverClickHouse = "clickhouse"
)

type Config struct {
	Driver                    string   `json:",default=mysql"`                                 // 数据库类型
	DSN                       string   `json:"dsn"`                                            // 数据库连接字符串
	MaxIdleConn               int      `json:",optional"`                                      // 最大空闲连接数
	MaxOpenConn               int      `json:",optional"`                                      // 最大连接数
	MaxLifetime               int      `json:",optional"`                                      // 连接最大生命周期(s)
	MaxIdleTime               int      `json:",optional"`                                      // 连接最大空闲时间(s)
	LogLevel                  string   `json:",default=info,options=[silent|error|warn|info]"` // 日志级别
	SlowThreshold             int      `json:",default=200"`                                   // 慢查询阈值(ms)
	Sources                   []string `json:",optional"`                                      // Master
	Replicas                  []string `json:",optional"`                                      // Slave
	SkipDefaultTransaction    bool     `json:",optional"`                                      // 跳过默认事务
	DefaultTransactionTimeout int      `json:",optional"`                                      // 默认事务超时时间(s) 	// 默认的 context.Context 超时时间
	DryRun                    bool     `json:",optional"`                                      // 启用“假执行”，只生成 SQL，不执行
	PrepareStmt               bool     `json:",optional"`                                      // 开启预编译语句缓存
	PrepareStmtMaxSize        int      `json:",optional"`                                      // 预编译语句缓存最大数量
	PrepareStmtTTL            int      `json:",optional"`                                      // 预编译语句缓存最大时间(s)
	DisableAutomaticPing      bool     `json:",optional"`                                      // 是否禁止 GORM 自动 ping 数据库以检测连接
	TranslateError            bool     `json:",default=true"`                                  // 是否转换错误
}

type TenantConfig struct {
	Default Config
	Tenants map[string]Config `json:",optional"`
}
