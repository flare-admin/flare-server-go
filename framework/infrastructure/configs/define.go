package configs

type Bootstrap struct {
	Server     *Server        `mapstructure:"server"`
	SuperAdmin *SuperAdmin    `mapstructure:"super_admin"`
	Log        *Log           `mapstructure:"log"`
	JWT        *JWT           `mapstructure:"jwt"`
	Data       *Data          `mapstructure:"data"`
	ConfPath   string         `mapstructure:"conf_path"`
	Storage    *StorageConfig `mapstructure:"storage"` // 添加存储配置
	NSQConfig  *NSQConfig     `mapstructure:"nsq"`
	NATSConfig *NATSConfig    `mapstructure:"nats"` // 添加 NATS 配置
}

type Server struct {
	Port               int    `mapstructure:"port"`
	RateQPS            int    `mapstructure:"rate_qps"`
	TracerPort         int    `mapstructure:"tracer_port"`
	Name               string `mapstructure:"name"`
	MaxRequestBodySize int    `mapstructure:"max_request_body_size"`
	TimeZone           string `mapstructure:"time_zone"`
}

type Log struct {
	OutPath    string `mapstructure:"output_dir"`
	FilePrefix string `mapstructure:"file_prefix"`
	Level      int64  `mapstructure:"level"`
	MaxSize    int64  `mapstructure:"max_size"`
	MaxBackups int64  `mapstructure:"max_backups"`
	MaxAge     int64  `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type JWT struct {
	Issuer            string `mapstructure:"issuer"`
	SigningKey        string `mapstructure:"signing_key"`
	ExpirationToken   int64  `mapstructure:"expiration_token"`
	ExpirationRefresh int64  `mapstructure:"expiration_refresh"`
}
type Data struct {
	DataBase *DataBase `mapstructure:"database"`
	Redis    *Redis    `mapstructure:"redis"`
}

// DataBase 数据库
type DataBase struct {
	EnableMigrate bool   `mapstructure:"enable_migrate"` // 是否开启迁移
	Driver        string `mapstructure:"driver"`
	Source        string `mapstructure:"source"`
	MaxIdleConns  int32  `mapstructure:"max_idle_conns"`
	MaxOpenConns  int32  `mapstructure:"max_open_conns"`
	LogLevel      int64  `mapstructure:"log_level"`
}

// Redis 数据库
type Redis struct {
	Network      string `mapstructure:"network"`
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	Db           int64  `mapstructure:"db"`
	ReadTimeout  int64  `mapstructure:"read_timeout"`
	WriteTimeout int64  `mapstructure:"write_timeout"`
}
type SuperAdmin struct {
	Nickname string `mapstructure:"nickname"`
	Phone    string `mapstructure:"phone"`
	Password string `mapstructure:"password"`
}
type NSQConfig struct {
	Address             string  `mapstructure:"address"`               // NSQ 地址
	MaxRetries          int     `mapstructure:"max_retries"`           // 最大重试次数
	MaxInFlight         int     `mapstructure:"max_in_flight"`         // 最大飞行消息数
	LookupdPollInterval int     `mapstructure:"lookupd_poll_interval"` // 查询间隔(毫秒)
	LookupdPollJitter   float64 `mapstructure:"lookupd_poll_jitter"`   // 查询抖动
	MaxBackoffDuration  int     `mapstructure:"max_backoff_duration"`  // 最大退避时间(毫秒)
	BackoffMultiplier   float64 `mapstructure:"backoff_multiplier"`    // 退避乘数
	Deflate             bool    `mapstructure:"deflate"`               // 是否启用压缩
	DeflateLevel        int     `mapstructure:"deflate_level"`         // 压缩级别
	Snappy              bool    `mapstructure:"snappy"`                // 是否启用 Snappy 压缩
}

// NATSConfig NATS 配置
type NATSConfig struct {
	// Address 地址
	Address string `mapstructure:"address"`
	// MaxRetries 最大重试次数
	MaxRetries int `mapstructure:"max_retries"`
	// ReconnectWait 重连等待时间(秒)
	ReconnectWait int `mapstructure:"reconnect_wait"`
	// MaxReconnects 最大重连次数
	MaxReconnects int `mapstructure:"max_reconnects"`
	// QueueGroup 队列组名称
	QueueGroup string `mapstructure:"queue_group"`
	// DurableName 持久化名称
	DurableName string `mapstructure:"durable_name"`
	// AckWait 确认等待时间(秒)
	AckWait int `mapstructure:"ack_wait"`
	// MaxDeliver 最大投递次数
	MaxDeliver int `mapstructure:"max_deliver"`
}
