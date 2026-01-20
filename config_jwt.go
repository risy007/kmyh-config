package config

type JWTConfig struct {
	SigningKey string `mapstructure:"SigningKey"` // JWT签名密钥
}
