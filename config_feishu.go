package config

import "fmt"

// FeishuConfig 飞书机器人配置
type FeishuConfig struct {
	Enabled           bool   `mapstructure:"enabled"`
	AppID             string `mapstructure:"app_id"`
	AppSecret         string `mapstructure:"app_secret"`
	EncryptKey        string `mapstructure:"encrypt_key"`
	VerificationToken string `mapstructure:"verification_token"`
	WebhookPath       string `mapstructure:"webhook_path"`

	// 会话共享配置
	SessionShared bool   `mapstructure:"session_shared"`
	SessionPrefix string `mapstructure:"session_prefix"`

	// OpenClaw 集成配置
	EnableOpenClaw bool   `mapstructure:"enable_openclaw"`
	OpenClawBotID  string `mapstructure:"openclaw_bot_id"`

	// NATS Micro 配置
	EnableMicro      bool   `mapstructure:"enable_micro"`
	MicroServiceName string `mapstructure:"micro_service_name"`
	MicroVersion     string `mapstructure:"micro_version"`
	MicroDescription string `mapstructure:"micro_description"`

	// 多用户配置
	Users []FeishuUserConfig `mapstructure:"users"`

	// 默认响应配置
	AckMessage     string `mapstructure:"ack_message"`
	TimeoutMessage string `mapstructure:"timeout_message"`
	UnknownCmdMsg  string `mapstructure:"unknown_cmd_msg"`
}

// Validate 验证飞书配置
func (cfg *FeishuConfig) Validate() error {
	if !cfg.Enabled {
		return nil
	}

	// 检查是否有至少一个机器人配置
	if cfg.AppID == "" && len(cfg.Users) == 0 {
		return fmt.Errorf("feishu app_id or users config is required when enabled")
	}

	// 验证用户配置
	for i, user := range cfg.Users {
		if err := user.Validate(); err != nil {
			return fmt.Errorf("feishu user[%d] config invalid: %w", i, err)
		}
	}

	return nil
}

// GetUserConfig 获取用户配置
func (cfg *FeishuConfig) GetUserConfig(userID string) (FeishuUserConfig, bool) {
	for _, user := range cfg.Users {
		if user.UserID == userID {
			return user, true
		}
	}
	return FeishuUserConfig{}, false
}

// GetOpenClawBotID 获取用户对应的OpenClaw Bot ID
func (cfg *FeishuConfig) GetOpenClawBotID(userID string) string {
	// 先查找用户特定配置
	for _, user := range cfg.Users {
		if user.UserID == userID && user.OpenClawBotID != "" {
			return user.OpenClawBotID
		}
	}
	// 使用全局配置
	return cfg.OpenClawBotID
}

// FeishuUserConfig 用户级别的飞书机器人配置
type FeishuUserConfig struct {
	UserID            string `mapstructure:"user_id"`
	AppID             string `mapstructure:"app_id"`
	AppSecret         string `mapstructure:"app_secret"`
	EncryptKey        string `mapstructure:"encrypt_key"`
	VerificationToken string `mapstructure:"verification_token"`
	Enabled           bool   `mapstructure:"enabled"`

	// 关联的 OpenClaw Bot
	OpenClawBotID string `mapstructure:"openclaw_bot_id"`

	// 用户特定的响应配置
	AckMessage   string `mapstructure:"ack_message"`
	SystemPrompt string `mapstructure:"system_prompt"`
}

// Validate 验证用户配置
func (cfg *FeishuUserConfig) Validate() error {
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
