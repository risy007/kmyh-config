package config

import "fmt"

// WeixinConfig 微信企业号配置
type WeixinConfig struct {
	Enabled           bool                `mapstructure:"enabled"`
	CorpID            string              `mapstructure:"corp_id"`
	WebHook           WorkwxWebHookConfig `mapstructure:"web_hook"`
	App               WorkwxAppConfig     `mapstructure:"app"`
	QYAPIHostOverride string              `mapstructure:"qyapi_host_override"`
	TLSKeyLogFile     string              `mapstructure:"tls_key_log_file"`
}

// Validate 验证微信配置
func (cfg *WeixinConfig) Validate() error {
	if cfg.Enabled && cfg.CorpID == "" {
		return fmt.Errorf("weixin corp ID is required when enabled")
	}
	return nil
}

// WorkwxWebHookConfig 企业微信WebHook配置
type WorkwxWebHookConfig struct {
	Key       string `mapstructure:"key"`
	Subscribe string `mapstructure:"subject"`
}

// Validate 验证企业微信WebHook配置
func (cfg *WorkwxWebHookConfig) Validate() error {
	if cfg.Key == "" {
		return fmt.Errorf("workwx webhook key is required")
	}
	return nil
}

// WorkwxAppConfig 企业微信应用配置
type WorkwxAppConfig struct {
	Address        string `mapstructure:"address"`
	CorpSecret     string `mapstructure:"corp_secret"`
	AgentID        int64  `mapstructure:"agent_id"`
	Token          string `mapstructure:"token"`
	EncodingAESKey string `mapstructure:"encoding_aes_key"`
	TxSubscribe    string `mapstructure:"tx_subject"`
	RxSubscribe    string `mapstructure:"rx_subject"`
}

// Validate 验证企业微信应用配置
func (cfg *WorkwxAppConfig) Validate() error {
	if cfg.Address == "" {
		return fmt.Errorf("workwx app address is required")
	}
	if cfg.CorpSecret == "" {
		return fmt.Errorf("workwx app corp secret is required")
	}
	if cfg.AgentID <= 0 {
		return fmt.Errorf("workwx app agent ID must be greater than 0")
	}
	if cfg.Token == "" {
		return fmt.Errorf("workwx app token is required")
	}
	return nil
}
