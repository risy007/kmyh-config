package config

import (
	"github.com/risy007/kmyh-config/local"
	"github.com/spf13/viper"
)

// AppConfig 应用主配置
type AppConfig struct {
	AppName string            `mapstructure:"name"`
	Env     string            `mapstructure:"env"`
	Etcd    config.EtcdConfig `mapstructure:"etcd"`
	Logger  config.LogConfig  `mapstructure:"logger"`
}

// LoadConfig 加载主配置文件
func NewAppConfig() (*AppConfig, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
