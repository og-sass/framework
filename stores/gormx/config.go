package gormx

const (
	DriverMysql      = "mysql"
	DriverPostgres   = "postgres"
	DriverClickHouse = "clickhouse"
)

type Config struct {
	Driver                    string   `json:"driver,default=mysql"`                                    // 数据库类型
	DSN                       string   `json:"dsn"`                                                     // 数据库连接字符串
	MaxIdleConn               int      `json:"max_idle_conn,optional"`                                  // 最大空闲连接数
	MaxOpenConn               int      `json:"max_open_conn,optional"`                                  // 最大连接数
	MaxLifetime               int      `json:"max_lifetime,optional"`                                   // 连接最大生命周期(s)
	MaxIdleTime               int      `json:"max_idle_time,optional"`                                  // 连接最大空闲时间(s)
	LogLevel                  string   `json:"log_level,default=info,options=[silent,error,warn,info]"` // 日志级别
	SlowThreshold             int      `json:"slow_threshold,default=200"`                              // 慢查询阈值(ms)
	Sources                   []string `json:"sources,optional"`                                        // Master
	Replicas                  []string `json:"replicas,optional"`                                       // Slave
	SkipDefaultTransaction    bool     `json:"skip_default_transaction,optional"`                       // 跳过默认事务
	DefaultTransactionTimeout int      `json:"default_transaction_timeout,optional"`                    // 默认事务超时时间(s)
	PrepareStmt               bool     `json:"prepare_stmt,optional"`                                   // 开启预编译语句缓存
	PrepareStmtMaxSize        int      `json:"prepare_stmt_max_size,optional"`                          // 预编译语句缓存最大数量
	PrepareStmtTTL            int      `json:"prepare_stmt_ttl,optional"`                               // 预编译语句缓存最大时间(s)
	DisableAutomaticPing      bool     `json:"disable_automatic_ping,optional"`                         // 是否禁止 GORM 自动 ping 数据库以检测连接
	TranslateError            bool     `json:"translate_error,default=true"`                            // 是否转换错误
	TenantDBName              string   `json:"tenant_db_name,default=site_id"`
}

type TenantConfig struct {
	Default      Config            `json:"default"`
	Tenants      map[string]Config `json:"tenants,optional"`
	TenantDBName string            `json:"tenant_db_name,default=site_id"`
}
