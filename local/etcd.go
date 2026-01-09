package config

import (
	"time"
)

// EtcdConfig etcd 连接配置
type EtcdConfig struct {
	Endpoints   []string      `yaml:"endpoints"`
	Username    string        `yaml:"username,omitempty"`
	Password    string        `yaml:"password,omitempty"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	TLS         *TLSConfig    `yaml:"tls,omitempty"`
	Prefix      string        `yaml:"prefix"`
}

// TLSConfig TLS 配置
type TLSConfig struct {
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
	CAFile   string `yaml:"ca_file"`
}
