package config

import "fmt"

// GoclawAdminConfig GoClaw 管理 API 配置
// 用于创建租户等管理操作
type GoclawAdminConfig struct {
	BaseURL string `mapstructure:"base_url"` // GoClaw API 地址 (http://host:port)
	APIKey  string `mapstructure:"api_key"`  // GoClaw Owner API Token
}

func (cfg *GoclawAdminConfig) Validate() error {
	if cfg.BaseURL == "" {
		return fmt.Errorf("goclaw admin base URL is required")
	}
	return nil
}

type GoclawBotConfig struct {
	UserID  string `mapstructure:"user_id"`  // 用户唯一标识，用于路由
	BaseURL string `mapstructure:"base_url"` // GoClaw WebSocket 地址 (ws://host:port)
	APIKey  string `mapstructure:"api_key"`  // GoClaw Token
	AgentID string `mapstructure:"agent_id"` // 目标 Agent ID
}

type GoclawConfig struct {
	List []GoclawBotConfig `mapstructure:"list"`

	Model           string `mapstructure:"model"`
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
}

func (cfg *GoclawConfig) Validate() error {
	if len(cfg.List) == 0 {
		return fmt.Errorf("goclaw bot list is empty")
	}
	hasDefault := false
	for _, bot := range cfg.List {
		if bot.UserID == "default" {
			hasDefault = true
		}
		if bot.BaseURL == "" {
			return fmt.Errorf("goclaw bot %s: base URL is required", bot.UserID)
		}
		if bot.APIKey == "" {
			return fmt.Errorf("goclaw bot %s: API key is required", bot.UserID)
		}
	}
	if !hasDefault {
		return fmt.Errorf("goclaw default bot is required in the list")
	}
	return nil
}
