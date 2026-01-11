package config

import (
	"fmt"
	"time"
)

// EtcdConfig 包含etcd连接的所有配置参数
// 用于分布式系统中配置的动态管理和更新
type EtcdConfig struct {
	// Endpoints etcd集群节点地址列表
	Endpoints []string `yaml:"endpoints" mapstructure:"endpoints"`
	// Username 认证用户名
	Username string `yaml:"username,omitempty" mapstructure:"username"`
	// Password 认证密码
	Password string `yaml:"password,omitempty" mapstructure:"password"`
	// DialTimeout 连接超时时间
	DialTimeout time.Duration `yaml:"dial_timeout" mapstructure:"dial_timeout"`
	// TLS TLS安全连接配置
	TLS *TLSConfig `yaml:"tls,omitempty" mapstructure:"tls"`
	// Prefix 配置键的前缀
	Prefix string `yaml:"prefix" mapstructure:"prefix"`
}

// Validate 验证etcd配置
func (cfg *EtcdConfig) Validate() error {
	if len(cfg.Endpoints) == 0 {
		return fmt.Errorf("etcd endpoints are required")
	}
	for _, endpoint := range cfg.Endpoints {
		if endpoint == "" {
			return fmt.Errorf("etcd endpoint cannot be empty")
		}
	}
	if cfg.DialTimeout <= 0 {
		return fmt.Errorf("dial timeout must be greater than 0")
	}
	return nil
}

// TLSConfig 包含TLS安全连接的配置参数
type TLSConfig struct {
	// CertFile 证书文件路径
	CertFile string `yaml:"cert_file" mapstructure:"cert_file"`
	// KeyFile 私钥文件路径
	KeyFile string `yaml:"key_file" mapstructure:"key_file"`
	// CAFile CA证书文件路径
	CAFile string `yaml:"ca_file" mapstructure:"ca_file"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level       string `mapstructure:"level"`
	Format      string `mapstructure:"format"`
	ToFile      bool   `mapstructure:"to_file"`
	Directory   string `mapstructure:"directory"`
	Development bool   `mapstructure:"development"`
}

// AppConfig 应用主配置
type AppConfig struct {
	AppName string     `mapstructure:"name"`
	Env     string     `mapstructure:"env"`
	Etcd    EtcdConfig `mapstructure:"etcd"`
	Logger  LogConfig  `mapstructure:"logger"`
}

// Validate 验证配置的有效性
func (cfg *AppConfig) Validate() error {
	if cfg.AppName == "" {
		return fmt.Errorf("app name is required")
	}
	if cfg.Env == "" {
		return fmt.Errorf("environment is required")
	}
	if len(cfg.Etcd.Endpoints) == 0 {
		return fmt.Errorf("etcd endpoints are required")
	}
	return nil
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Engine      string `mapstructure:"engine"`
	Name        string `mapstructure:"name"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	TablePrefix string `mapstructure:"table_prefix"`
	Parameters  string `mapstructure:"parameters"`

	MaxLifetime  int `mapstructure:"max_lifetime"`
	MaxOpenConns int `mapstructure:"max_open_conns"`
	MaxIdleConns int `mapstructure:"max_idle_conns"`
}

// Validate 验证数据库配置
func (cfg *DatabaseConfig) Validate() error {
	if cfg.Engine == "" {
		return fmt.Errorf("database engine is required")
	}
	if cfg.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if cfg.Port <= 0 {
		return fmt.Errorf("database port must be greater than 0")
	}
	if cfg.Name == "" {
		return fmt.Errorf("database name is required")
	}
	return nil
}

// HttpConfig HTTP服务配置
type HttpConfig struct {
	Listen string `mapstructure:"listen"`
	Prefix string `mapstructure:"prefix"`
}

// Validate 验证HTTP配置
func (cfg *HttpConfig) Validate() error {
	if cfg.Listen == "" {
		return fmt.Errorf("http listen address is required")
	}
	return nil
}

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

// FuiouConfig 富友支付配置
type FuiouConfig struct {
	MchntKey string `mapstructure:"mchnt_key"`
}

// Validate 验证富友支付配置
func (cfg *FuiouConfig) Validate() error {
	if cfg.MchntKey == "" {
		return fmt.Errorf("fuiou merchant key is required")
	}
	return nil
}

// NatsConfig NATS消息队列配置
type NatsConfig struct {
	Address    string   `mapstructure:"address"`
	Username   string   `mapstructure:"username"`
	Password   string   `mapstructure:"password"`
	Subscribes []string `mapstructure:"subscribes"`
}

// Validate 验证NATS配置
func (cfg *NatsConfig) Validate() error {
	if cfg.Address == "" {
		return fmt.Errorf("nats address is required")
	}
	return nil
}

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

// MiddleConfig 中间件配置
type MiddleConfig struct {
	IPWhiteList IpWhiteListConfig `mapstructure:"ip_whitelist"`
}

// Validate 验证中间件配置
func (cfg *MiddleConfig) Validate() error {
	return cfg.IPWhiteList.Validate()
}

// IpWhiteListConfig IP白名单配置
type IpWhiteListConfig struct {
	Enabled   bool     `mapstructure:"enabled"`
	WhiteList []string `mapstructure:"white_list"`
}

// Validate 验证IP白名单配置
func (cfg *IpWhiteListConfig) Validate() error {
	if cfg.Enabled && len(cfg.WhiteList) == 0 {
		return fmt.Errorf("IP whitelist is enabled but no IPs are provided")
	}
	return nil
}
