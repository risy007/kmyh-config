package config

type EmailConfig struct {
	SMTPHost              string  `mapstructure:"smtp_host"`
	SMTPPort              int     `mapstructure:"smtp_port"`
	Username              string  `mapstructure:"username"`
	Password              string  `mapstructure:"password"`
	FromAddress           string  `mapstructure:"from_address"`
	UseTLS                bool    `mapstructure:"use_tls"`
	TLSInsecureSkipVerify bool    `mapstructure:"tls_insecure_skip_verify"`
	RateLimitPerSecond    float64 `mapstructure:"rate_limit_per_second"`
	RateLimitBurst        int     `mapstructure:"rate_limit_burst"`
}
