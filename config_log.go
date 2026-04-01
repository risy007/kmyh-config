package config

// LogConfig 日志配置
type LogConfig struct {
	Level       string `mapstructure:"level" description:"日志级别"`
	Format      string `mapstructure:"format" description:"日志格式"`
	ToFile      bool   `mapstructure:"to_file" description:"是否写入文件"`
	Directory   string `mapstructure:"directory" description:"日志目录"`
	Development bool   `mapstructure:"development" description:"开发模式"`
}
