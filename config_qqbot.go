package config

import "fmt"

// QQBotUserConfig QQ机器人用户配置
type QQBotUserConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	UserID    string `mapstructure:"user_id"`
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
}

// QQBotConfig QQ机器人配置
type QQBotConfig struct {
	Enabled        bool              `mapstructure:"enabled"`
	AppID          string            `mapstructure:"app_id"`
	AppSecret      string            `mapstructure:"app_secret"`
	Token          string            `mapstructure:"token"`
	WebhookPath    string            `mapstructure:"webhook_path"`
	UseSandbox     bool              `mapstructure:"use_sandbox"`
	SandboxGuildID string            `mapstructure:"sandbox_guild_id"`
	Users          []QQBotUserConfig `mapstructure:"users"`

	// 会话共享配置
	SessionShared bool   `mapstructure:"session_shared"`
	SessionPrefix string `mapstructure:"session_prefix"`

	// OpenClaw 集成配置 (通过 ChannelMapper 统一管理)
	EnableOpenClaw bool `mapstructure:"enable_openclaw"`

	// NATS Micro 配置
	EnableMicro      bool   `mapstructure:"enable_micro"`
	MicroServiceName string `mapstructure:"micro_service_name"`
	MicroVersion     string `mapstructure:"micro_version"`
	MicroDescription string `mapstructure:"micro_description"`

	AckMessage     string `mapstructure:"ack_message"`
	TimeoutMessage string `mapstructure:"timeout_message"`
	UnknownCmdMsg  string `mapstructure:"unknown_cmd_msg"`
}

// Validate 验证QQBot配置
func (cfg *QQBotConfig) Validate() error {
	if !cfg.Enabled {
		return nil
	}

	// 检查是否有至少一个机器人配置
	if cfg.AppID == "" && len(cfg.Users) == 0 {
		return fmt.Errorf("qqbot app_id or users config is required when enabled")
	}

	// 验证用户配置
	for i, user := range cfg.Users {
		if err := user.Validate(); err != nil {
			return fmt.Errorf("qqbot user[%d] config invalid: %w", i, err)
		}
	}

	return nil
}

// GetUserConfig 获取用户配置
func (cfg *QQBotConfig) GetUserConfig(userID string) (QQBotUserConfig, bool) {
	for _, user := range cfg.Users {
		if user.UserID == userID {
			return user, true
		}
	}
	return QQBotUserConfig{}, false
}

// Validate 验证用户配置
func (cfg *QQBotUserConfig) Validate() error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if cfg.AppID == "" {
		return fmt.Errorf("app_id is required")
	}
	if cfg.AppSecret == "" {
		return fmt.Errorf("app_secret is required")
	}
	return nil
}
