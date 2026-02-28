package config

import "fmt"

// WeixinConfig 微信企业号配置
type WeixinConfig struct {
	Enabled           bool                `mapstructure:"enabled"`
	CorpID            string              `mapstructure:"corp_id"`
	WebHook           WorkwxWebHookConfig `mapstructure:"web_hook"`
	App               WorkwxAppConfig     `mapstructure:"app"`
	Event             WeixinEventConfig   `mapstructure:"event"`
	Micro             WeixinMicroConfig   `mapstructure:"micro"`
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
	PublishRx      bool   `mapstructure:"publish_rx"`
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

// WeixinEventConfig 企业微信事件配置
type WeixinEventConfig struct {
	EnableContactChange   bool   `mapstructure:"enable_contact_change"`   // 通讯录变更事件
	EnableExternalContact bool   `mapstructure:"enable_external_contact"` // 外部联系人事件
	EnableOAApproval      bool   `mapstructure:"enable_oa_approval"`      // OA审批事件
	EnableAppMenu         bool   `mapstructure:"enable_app_menu"`         // 应用菜单事件
	EnableAppSubscribe    bool   `mapstructure:"enable_app_subscribe"`    // 关注/取消关注事件
	EnableExternalChat    bool   `mapstructure:"enable_external_chat"`    // 客户群变更事件
	EnableKF              bool   `mapstructure:"enable_kf"`               // 客服消息事件
	PublishEvents         bool   `mapstructure:"publish_events"`          // 发布事件到NATS
	EventSubject          string `mapstructure:"event_subject"`           // 事件发布主题
}

// WeixinMicroConfig 企业微信 NATS Micro 服务配置
type WeixinMicroConfig struct {
	Enable         bool   `mapstructure:"enable"`          // 是否启用 Micro 服务
	ServiceName    string `mapstructure:"service_name"`    // 服务名称
	Version        string `mapstructure:"version"`         // 服务版本
	Description    string `mapstructure:"description"`     // 服务描述
	EnableWebhook  bool   `mapstructure:"enable_webhook"`  // 启用 Webhook 消息
	EnableAppMsg   bool   `mapstructure:"enable_app_msg"`  // 启用 App 消息
	EnableContact  bool   `mapstructure:"enable_contact"`  // 启用通讯录 API
	EnableChat     bool   `mapstructure:"enable_chat"`     // 启用群聊 API
	EnableExternal bool   `mapstructure:"enable_external"` // 启用外部联系人 API
	EnableMedia    bool   `mapstructure:"enable_media"`    // 启用媒体文件 API
	EnableAgent    bool   `mapstructure:"enable_agent"`    // 启用应用管理 API
	EnableUser     bool   `mapstructure:"enable_user"`     // 启用用户身份 API
}
