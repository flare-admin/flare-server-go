package log

type LogConf struct {
	OutPath    string `json:"out_path" yaml:"out_path"`
	FilePrefix string `json:"file_prefix"`
	Level      int64  `json:"level"`
	MaxSize    int64  `json:"max_size"`
	MaxBackups int64  `json:"max_backups"`
	MaxAge     int64  `json:"max_age"`
	Compress   bool   `json:"compress"`
}

// NewLogConf creates a new LogConf instance with the provided options
func NewLogConf(options ...Option) *LogConf {
	// Create a default configuration
	conf := NewDef()

	// Apply each option
	for _, option := range options {
		option(conf)
	}

	return conf
}

// Option is a function that configures a LogConf
type Option func(*LogConf)

// NewDef creates a default LogConf
func NewDef() *LogConf {
	return &LogConf{
		OutPath:    "./logs",
		FilePrefix: "ares",
		Level:      1,
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	}
}

// WithOutPath sets the output path for the log files
func WithOutPath(outPath string) Option {
	return func(c *LogConf) {
		c.OutPath = outPath
	}
}

// WithFilePrefix sets the file prefix for log files
func WithFilePrefix(prefix string) Option {
	return func(c *LogConf) {
		c.FilePrefix = prefix
	}
}

// WithLevel sets the log level
func WithLevel(level int64) Option {
	return func(c *LogConf) {
		c.Level = level
	}
}

// WithMaxSize sets the maximum size of each log file
func WithMaxSize(size int64) Option {
	return func(c *LogConf) {
		c.MaxSize = size
	}
}

// WithMaxBackups sets the maximum number of log backups
func WithMaxBackups(backups int64) Option {
	return func(c *LogConf) {
		c.MaxBackups = backups
	}
}

// WithMaxAge sets the maximum age of log files (in days)
func WithMaxAge(age int64) Option {
	return func(c *LogConf) {
		c.MaxAge = age
	}
}

// WithCompress enables or disables log file compression
func WithCompress(compress bool) Option {
	return func(c *LogConf) {
		c.Compress = compress
	}
}
