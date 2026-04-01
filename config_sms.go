package config

type AliyunSMSConfig struct {
	AccessKeyID        string  `mapstructure:"access_key_id" description:"AccessKey ID"`
	AccessKeySecret    string  `mapstructure:"access_key_secret" description:"AccessKey 密钥"`
	RegionID           string  `mapstructure:"region_id" description:"区域ID"`
	SignName           string  `mapstructure:"sign_name" description:"短信签名"`
	HTTPTimeout        int64   `mapstructure:"http_timeout" description:"HTTP超时时间(秒)"`
	RateLimitPerSecond float64 `mapstructure:"rate_limit_per_second" description:"每秒速率限制"`
	RateLimitBurst     int     `mapstructure:"rate_limit_burst" description:"突发请求数限制"`
}
