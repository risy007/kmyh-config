package config

type AuthConfig struct {
	Enable             bool     `mapstructure:"enable"`
	TokenExpired        string   `mapstructure:"token_expired"`         // time.ParseDuration(),示例：1d,10m,50s....
	RefreshTokenExpired string   `mapstructure:"refresh_token_expired"` // 刷新令牌过期时间
	IgnorePathPrefixes  []string `mapstructure:"ignore_path_prefixes"`
	JWTSigningKey      string   `mapstructure:"jwt_signing_key"`
	Issuer             string   `mapstructure:"issuer"`
	VerifyModes        []string `mapstructure:"verify_modes"` // 验证模式: captcha, sms, email
	CodeExpired        string   `mapstructure:"code_expired"` // 验证码过期时间，time.ParseDuration()格式，如: "120s", "5m", "1h"
	RedisPrefix        string   `mapstructure:"redis_prefix"`
	KeyLong            int      `mapstructure:"key_long"`
	KeyContent         string   `mapstructure:"key_content"` // numeric, alphabetic, alphanumeric
}
