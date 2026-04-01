package config

import "fmt"

type DifyBotConfig struct {
	UserID     string `mapstructure:"user_id"`
	BaseURL    string `mapstructure:"base_url"`
	APIKey     string `mapstructure:"api_key"`
	WorkflowID string `mapstructure:"workflow_id"`
	BotType    string `mapstructure:"bot_type"`
}

type DifyConfig struct {
	List []DifyBotConfig `mapstructure:"list"`

	Model string `mapstructure:"model"`
	// BaseURL 和 APIKey 可选，如果设置了就是默认值
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`

	CachePeriod     string `mapstructure:"cache_period"`
	HandleWorkwxMsg bool   `mapstructure:"handle_workwx_msg"`
	NotifyUsers     string `mapstructure:"notify_users"`

	Timeout            string `mapstructure:"timeout"`
	StreamOutput       bool   `mapstructure:"stream_output"`
	StreamChunkTimeout string `mapstructure:"stream_chunk_timeout"`
	MaxRetries         int    `mapstructure:"max_retries"`
	TotalTimeout       string `mapstructure:"total_timeout"`

	AckMessage     string `mapstructure:"ack_message"`
	TimeoutMessage string `mapstructure:"timeout_message"`

	DefaultPrompt string `mapstructure:"default_prompt"`
}

// GetEffectiveBaseURL 获取有效的 BaseURL
func (cfg *DifyBotConfig) GetEffectiveBaseURL(defaultURL string) string {
	if cfg.BaseURL != "" {
		return cfg.BaseURL
	}
	return defaultURL
}

// GetEffectiveAPIKey 获取有效的 APIKey
func (cfg *DifyBotConfig) GetEffectiveAPIKey(defaultKey string) string {
	if cfg.APIKey != "" {
		return cfg.APIKey
	}
	return defaultKey
}

func (cfg *DifyConfig) Validate() error {
	if len(cfg.List) == 0 {
		return fmt.Errorf("dify bot list is empty")
	}
	hasDefault := false
	for _, bot := range cfg.List {
		if bot.UserID == "default" {
			hasDefault = true
		}
		if bot.BaseURL == "" {
			return fmt.Errorf("dify bot %s: base URL is required", bot.UserID)
		}
		if bot.APIKey == "" {
			return fmt.Errorf("dify bot %s: API key is required", bot.UserID)
		}
	}
	if !hasDefault {
		return fmt.Errorf("dify default bot is required in the list")
	}
	return nil
}
