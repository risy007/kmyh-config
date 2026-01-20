package config

type AliyunSMSConfig struct {
	AccessKeyID        string  `mapstructure:"access_key_id"`
	AccessKeySecret    string  `mapstructure:"access_key_secret"`
	RegionID           string  `mapstructure:"region_id"`
	SignName           string  `mapstructure:"sign_name"`
	HTTPTimeout        int64   `mapstructure:"http_timeout"`
	RateLimitPerSecond float64 `mapstructure:"rate_limit_per_second"`
	RateLimitBurst     int     `mapstructure:"rate_limit_burst"`
}
