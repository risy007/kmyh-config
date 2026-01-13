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
