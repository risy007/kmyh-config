package config

import "fmt"

type RedisConfig struct {
	Host      string `mapstructure:"host" description:"服务器地址"`
	Port      int    `mapstructure:"port" description:"端口号"`
	Password  string `mapstructure:"password" description:"密码"`
	KeyPrefix string `mapstructure:"key_prefix" description:"键前缀"`
	MainDBId  int    `mapstructure:"main_db_id" description:"主库ID"`
}

func (cfg *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}
