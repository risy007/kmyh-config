package config

import "fmt"

// FuiouConfig 富友支付配置
type FuiouConfig struct {
	MchntKey        string `mapstructure:"mchnt_key"`
	PublishToNATS   bool   `mapstructure:"publish_to_nats"`   // 是否发布到 NATS
	PublishToWeixin bool   `mapstructure:"publish_to_weixin"` // 是否直接发送到企业微信
}

// Validate 验证富友支付配置
func (cfg *FuiouConfig) Validate() error {
	if cfg.MchntKey == "" {
		return fmt.Errorf("fuiou merchant key is required")
	}
	return nil
}
