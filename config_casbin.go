package config

type CasbinConfig struct {
	Enable             bool     `mapstructure:"Enable"`
	Debug              bool     `mapstructure:"Debug"`
	Model              string   `mapstructure:"Model"`
	AutoLoad           bool     `mapstructure:"AutoLoad"`
	AutoLoadInternal   int      `mapstructure:"AutoLoadInternal"`
	IgnorePathPrefixes []string `mapstructure:"IgnorePathPrefixes"`
}
