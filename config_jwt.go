package config

type JWTConfig struct {
	SigningKey string `mapstructure:"signing_key" description:"JWT签名密钥"`
	Issuer     string `mapstructure:"issuer" description:"签发者"`
	Expire     int    `mapstructure:"expire" description:"过期时间(小时)"`
}
