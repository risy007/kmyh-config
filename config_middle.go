package config

import "fmt"

// MiddleConfig 中间件配置
type MiddleConfig struct {
	IPWhiteList IpWhiteListConfig `mapstructure:"ip_whitelist"`
}

// Validate 验证中间件配置
func (cfg *MiddleConfig) Validate() error {
	return cfg.IPWhiteList.Validate()
}

// IpWhiteListConfig IP白名单配置
type IpWhiteListConfig struct {
	Enabled   bool     `mapstructure:"enabled"`
	WhiteList []string `mapstructure:"white_list"`
}

// Validate 验证IP白名单配置
func (cfg *IpWhiteListConfig) Validate() error {
	if cfg.Enabled && len(cfg.WhiteList) == 0 {
		return fmt.Errorf("IP whitelist is enabled but no IPs are provided")
	}
	return nil
}
