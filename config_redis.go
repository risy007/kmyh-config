package config

import "fmt"

type RedisConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Password  string `mapstructure:"password"`
	KeyPrefix string `mapstructure:"key_prefix"`
	MainDBId  int    `mapstructure:"main_db_id"`
}

func (cfg *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}
