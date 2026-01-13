package config

type RedisConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Password  string `mapstructure:"password"`
	KeyPrefix string `mapstructure:"key_prefix"`
}
