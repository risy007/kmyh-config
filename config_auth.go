package config

type AuthConfig struct {
	Enable             bool     `mapstructure:"Enable"`
	TokenExpired       string   `mapstructure:"TokenExpired"` // time.ParseDuration(),示例：1d,10m,50s....
	IgnorePathPrefixes []string `mapstructure:"IgnorePathPrefixes"`
	JWTSigningKey      string   `mapstructure:"JWTSigningKey"`
	Issuer             string   `mapstructure:"Issuer"`
	VerifyModes        []string `mapstructure:"VerifyModes"` // 验证模式: captcha, sms, email
	CodeExpired        string   `mapstructure:"CodeExpired"` // 验证码过期时间，time.ParseDuration()格式，如: "120s", "5m", "1h"
	RedisPrefix        string   `mapstructure:"redis_prefix"`
	KeyLong            int      `mapstructure:"key_long"`
	KeyContent         string   `mapstructure:"key_content"` // numeric, alphabetic, alphanumeric
}
