package config

import "fmt"

// PrtgConfig PRTG网络监控配置
type PrtgConfig struct {
	Subject string `mapstructure:"mq_subject"`
}

// Validate 验证PRTG配置
func (cfg *PrtgConfig) Validate() error {
	if cfg.Subject == "" {
		return fmt.Errorf("prtg subject is required")
	}
	return nil
}
