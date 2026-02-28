package config

import "fmt"

type OpenclawConfig struct {
	BaseURL         string `mapstructure:"base_url"`
	APIKey          string `mapstructure:"api_key"`
	Model           string `mapstructure:"model"`
	CachePeriod     string `mapstructure:"cache_period"`
	HandleWorkwxMsg bool   `mapstructure:"handle_workwx_msg"`
	NotifyUsers     string `mapstructure:"notify_users"`
}

func (cfg *OpenclawConfig) Validate() error {
	if cfg.BaseURL == "" {
		return fmt.Errorf("openclaw base URL is required")
	}
	if cfg.APIKey == "" {
		return fmt.Errorf("openclaw API key is required")
	}
	return nil
}
