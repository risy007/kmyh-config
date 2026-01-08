package config

type NatsConfig struct {
	Address    string   `mapstructure:"address"`
	Username   string   `mapstructure:"username"`
	Password   string   `mapstructure:"password"`
	Subscribes []string `mapstructure:"subscribes"`
}
