package config

import "fmt"

// FuiouConfig 富友支付配置
type FuiouConfig struct {
	MchntKey string `mapstructure:"mchnt_key"`
}

// Validate 验证富友支付配置
func (cfg *FuiouConfig) Validate() error {
	if cfg.MchntKey == "" {
		return fmt.Errorf("fuiou merchant key is required")
	}
	return nil
}
