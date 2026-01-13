package config

// LogConfig 日志配置
type LogConfig struct {
	Level       string `mapstructure:"level"`
	Format      string `mapstructure:"format"`
	ToFile      bool   `mapstructure:"to_file"`
	Directory   string `mapstructure:"directory"`
	Development bool   `mapstructure:"development"`
}
