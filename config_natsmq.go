package config

import "fmt"

// NatsConfig NATS 连接配置（从 etcd 加载）
//
// etcd 配置示例:
//
//	/configs/{app}/{env}/nats/content.yaml
//	address: tls://172.18.0.194:4222
//	token: svc_xxxx
//	auto_detect_machine_id: true
//	tls:
//	  ca_file: /opt/kmyh/yh-admin/certs/nats-ca.pem
//	reconnect_wait: 2
//	max_reconnects: 60
type NatsConfig struct {
	Address string `mapstructure:"address"`

	Token               string `mapstructure:"token"`
	MachineID           string `mapstructure:"machine_id"`
	AutoDetectMachineID bool   `mapstructure:"auto_detect_machine_id"`

	TLS *TLSConfig `mapstructure:"tls"`

	ReconnectWait int `mapstructure:"reconnect_wait"`
	MaxReconnects int `mapstructure:"max_reconnects"`
}

func (cfg *NatsConfig) Validate() error {
	if cfg.Address == "" {
		return fmt.Errorf("nats address is required")
	}
	return nil
}

func (cfg *NatsConfig) UseTokenAuth() bool {
	return cfg.Token != ""
}
