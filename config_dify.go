package config

import "fmt"

// DifyConfig Dify AI平台配置
type DifyConfig struct {
	BaseURL       string `mapstructure:"base_url"`
	APIKey        string `mapstructure:"api_key"`
	CachePeriod   string `mapstructure:"cache_period"`
	DefaultPrompt string `mapstructure:"default_prompt"`
	BotType       string `mapstructure:"bot_type"`
	WorkflowID    string `mapstructure:"workflow_id"`
}

// Validate 验证Dify配置
func (cfg *DifyConfig) Validate() error {
	if cfg.BaseURL == "" {
		return fmt.Errorf("dify base URL is required")
	}
	if cfg.APIKey == "" {
		return fmt.Errorf("dify API key is required")
	}
	return nil
}
