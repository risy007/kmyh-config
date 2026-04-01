package config

type EmailConfig struct {
	SMTPHost              string  `mapstructure:"smtp_host" description:"SMTP服务器地址"`
	SMTPPort              int     `mapstructure:"smtp_port" description:"SMTP端口号"`
	Username              string  `mapstructure:"username" description:"用户名"`
	Password              string  `mapstructure:"password" description:"密码"`
	FromAddress           string  `mapstructure:"from_address" description:"发件人地址"`
	UseTLS                bool    `mapstructure:"use_tls" description:"使用TLS加密"`
	TLSInsecureSkipVerify bool    `mapstructure:"tls_insecure_skip_verify" description:"跳过TLS证书验证"`
	RateLimitPerSecond    float64 `mapstructure:"rate_limit_per_second" description:"每秒速率限制"`
	RateLimitBurst        int     `mapstructure:"rate_limit_burst" description:"突发请求数限制"`
}
