package config

import "fmt"

type OpenclawBotConfig struct {
	SystemUserID string `mapstructure:"system_user_id"` // 系统唯一用户ID，通过 ChannelMapper 关联
	BaseURL      string `mapstructure:"base_url"`
	APIKey       string `mapstructure:"api_key"`
	Model        string `mapstructure:"model"` // 可选，默认使用 OpenclawConfig.Model
}

type OpenclawConfig struct {
	List []OpenclawBotConfig `mapstructure:"list"`

	Model           string `mapstructure:"model"`             // 默认模型
	CachePeriod     string `mapstructure:"cache_period"`
	HandleWorkwxMsg bool   `mapstructure:"handle_workwx_msg"`
	NotifyUsers     string `mapstructure:"notify_users"`

	// 超时和流式输出配置
	Timeout            string `mapstructure:"timeout"`              // 单次请求超时时间（如 "60s", "1m"）
	StreamOutput       bool   `mapstructure:"stream_output"`        // 是否启用流式输出
	StreamChunkTimeout string `mapstructure:"stream_chunk_timeout"` // 流式输出的分块超时时间
	MaxRetries         int    `mapstructure:"max_retries"`          // 流式输出时的最大重试次数
	TotalTimeout       string `mapstructure:"total_timeout"`        // 流式输出的总超时时间

	// 消息配置
	AckMessage     string `mapstructure:"ack_message"`     // 收到消息后的确认消息内容
	TimeoutMessage string `mapstructure:"timeout_message"` // 超时后的通知消息模板
}

// GetEffectiveModel 获取有效的模型，如果 bot 没有指定则使用默认模型
func (cfg *OpenclawBotConfig) GetEffectiveModel(defaultModel string) string {
	if cfg.Model != "" {
		return cfg.Model
	}
	return defaultModel
}

func (cfg *OpenclawConfig) Validate() error {
	if len(cfg.List) == 0 {
		return fmt.Errorf("openclaw bot list is empty")
	}
	hasDefault := false
	for _, bot := range cfg.List {
		if bot.SystemUserID == "default" {
			hasDefault = true
		}
		if bot.BaseURL == "" {
			return fmt.Errorf("openclaw bot %s: base URL is required", bot.SystemUserID)
		}
		if bot.APIKey == "" {
			return fmt.Errorf("openclaw bot %s: API key is required", bot.SystemUserID)
		}
	}
	if !hasDefault {
		return fmt.Errorf("openclaw default bot is required in the list")
	}
	return nil
}
