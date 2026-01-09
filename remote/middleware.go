package config

type (
	IpWhiteListConfig struct {
		Enabled   bool     `mapstructure:"enabled"`
		WhiteList []string `mapstructure:"white_list"`
	}
	MiddleConfig struct {
		IPWhiteList IpWhiteListConfig `mapstructure:"IPwhitelist"`
	}
)
