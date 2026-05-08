package config

import (
	"fmt"
)

// HttpConfig HTTP服务配置
type HttpConfig struct {
	Host        string `mapstructure:"host" description:"监听地址"`
	Port        int    `mapstructure:"port" description:"端口号"`
	Prefix      string `mapstructure:"prefix" description:"路由前缀"`
	ExternalURL string `mapstructure:"external_url" description:"外部访问地址"`
}

// ListenAddr 返回监听地址
func (a *HttpConfig) ListenAddr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

// Validate 验证HTTP配置
func (a *HttpConfig) Validate() error {
	if a.Host == "" {
		return fmt.Errorf("http host is required")
	}
	if a.Port <= 0 {
		return fmt.Errorf("http port must be greater than 0")
	}
	return nil
}
