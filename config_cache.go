package config

type CacheConfig struct {
	Period string `mapstructure:"period" description:"缓存周期"`
}
