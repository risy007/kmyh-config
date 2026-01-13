package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// AppConfig 应用主配置
type AppConfig struct {
	AppName string     `mapstructure:"name"`
	Env     string     `mapstructure:"env"`
	Etcd    EtcdConfig `mapstructure:"etcd"`
	Logger  LogConfig  `mapstructure:"logger"`
}

// Validate 验证配置的有效性
func (cfg *AppConfig) Validate() error {
	if cfg.AppName == "" {
		return fmt.Errorf("app name is required")
	}
	if cfg.Env == "" {
		return fmt.Errorf("environment is required")
	}
	if len(cfg.Etcd.Endpoints) == 0 {
		return fmt.Errorf("etcd endpoints are required")
	}
	return nil
}

// NewAppConfig 从配置文件创建并加载主配置
// 默认会查找当前目录及config子目录下的config.yaml文件
// 返回配置对象和可能的错误
func NewAppConfig() (*AppConfig, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}
