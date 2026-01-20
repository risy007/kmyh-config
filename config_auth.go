package config

type AuthConfig struct {
	Enable             bool     `mapstructure:"Enable"`
	TokenExpired       string   `mapstructure:"TokenExpired"` //time.ParseDuration(),示例：1d,10m,50s.....
	IgnorePathPrefixes []string `mapstructure:"IgnorePathPrefixes"`
	JWTSigningKey      string   `mapstructure:"JWTSigningKey"`
	Issuer             string   `mapstructure:"Issuer"`
	VerifyModes        []string `mapstructure:"VerifyModes"`   // 验证模式: captcha, sms, email
	ExpireMinutes      int      `mapstructure:"ExpireMinutes"` // 验证码有效期（分钟）
}
