package config

type (
	HttpConfig struct {
		Listen string `mapstructure:"listen"`
		Prefix string `mapstructure:"prefix"`
	}
)
