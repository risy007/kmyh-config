package config

import "fmt"

// NatsConfig NATS消息队列配置
type NatsConfig struct {
	Address    string   `mapstructure:"address"`
	Username   string   `mapstructure:"username"`
	Password   string   `mapstructure:"password"`
	Subscribes []string `mapstructure:"subscribes"`
}

// Validate 验证NATS配置
func (cfg *NatsConfig) Validate() error {
	if cfg.Address == "" {
		return fmt.Errorf("nats address is required")
	}
	return nil
}
