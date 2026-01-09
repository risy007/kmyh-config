package config

type LogConfig struct {
	Level       string `mapstructure:"Level"`
	Format      string `mapstructure:"Format"`
	ToFile      bool   `mapstructure:"ToFile"`
	Directory   string `mapstructure:"Directory"`
	Development bool   `mapstructure:"Development"`
}
