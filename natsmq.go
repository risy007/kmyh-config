package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// NewNatsConnection 根据 etcd 中的 NatsConfig 创建 NATS 连接
// appName 用作客户端名称（显示在 NATS 服务端连接列表）
func NewNatsConnection(cfg NatsConfig, appName string, logger *zap.Logger) (*nats.Conn, error) {
	log := logger.With(zap.Namespace("[natsmq]")).Sugar()

	var options []nats.Option

	reconnectWait := 2 * time.Second
	if cfg.ReconnectWait > 0 {
		reconnectWait = time.Duration(cfg.ReconnectWait) * time.Second
	}
	maxReconnects := 60
	if cfg.MaxReconnects > 0 {
		maxReconnects = cfg.MaxReconnects
	}
	options = append(options,
		nats.ReconnectWait(reconnectWait),
		nats.MaxReconnects(maxReconnects),
	)

	options = append(options,
		nats.ConnectHandler(func(nc *nats.Conn) {
			log.Info("NATS 连接成功",
				zap.String("address", nc.ConnectedUrl()),
				zap.String("server_id", nc.ConnectedServerId()),
			)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Info("NATS 重连成功",
				zap.String("address", nc.ConnectedUrl()),
			)
		}),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				log.Warn("NATS 连接断开", zap.Error(err))
			}
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Error("NATS 连接关闭", zap.Error(nc.LastError()))
		}),
	nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
		if sub != nil {
			log.Error("NATS 异常", zap.Error(err), zap.String("subject", sub.Subject))
		} else {
			log.Error("NATS 异常", zap.Error(err))
		}
	}),
	)

	machineID := cfg.MachineID
	if cfg.AutoDetectMachineID || machineID == "" {
		id, err := machineid.ID()
		if err != nil {
			return nil, fmt.Errorf("failed to detect machine ID: %w", err)
		}
		machineID = id
		log.Info("自动检测 MachineID", zap.String("machine_id", machineID))
	}

	// 始终设置连接名，供 NATS 服务端和 auth callout 识别客户端
	if appName != "" {
		options = append(options, nats.Name(appName))
	}

	if cfg.UseTokenAuth() {
		options = append(options, nats.UserInfo(cfg.Token, machineID))
		log.Info("使用 Token+MachineID 认证",
			zap.String("token_prefix", cfg.Token[:min(12, len(cfg.Token))]),
			zap.String("machine_id", machineID),
		)
	} else {
		options = append(options, nats.UserInfo("", machineID))
		log.Info("无 Token，发送 MachineID 用于自动注册",
			zap.String("machine_id", machineID),
		)
	}

	needNatsTLS := cfg.TLS != nil || strings.HasPrefix(cfg.Address, "tls://")
	if needNatsTLS {
		tlsCfg := cfg.TLS
		if tlsCfg == nil {
			tlsCfg = &TLSConfig{}
		}
		tlsConfig, err := createNatsTLSConfig(tlsCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create NATS TLS config: %w", err)
		}
		options = append(options, nats.Secure(tlsConfig))
		log.Info("已启用 TLS 加密连接")
	}

	nc, err := nats.Connect(cfg.Address, options...)
	if err != nil {
		log.Error("NATS 连接失败", zap.Error(err))
		return nil, err
	}

	return nc, nil
}

// createNatsTLSConfig 根据 TLS 配置创建 tls.Config (用于 NATS)
func createNatsTLSConfig(cfg *TLSConfig) (*tls.Config, error) {
	caCert := embeddedCA
	if cfg.CAFile != "" {
		data, err := ioutil.ReadFile(cfg.CAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA file: %w", err)
		}
		caCert = data
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}

	tlsConfig := &tls.Config{
		RootCAs:    caCertPool,
		MinVersion: tls.VersionTLS12,
	}

	if cfg.CertFile != "" && cfg.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	} else if len(embeddedClientCert) > 0 && len(embeddedClientKey) > 0 {
		cert, err := tls.X509KeyPair(embeddedClientCert, embeddedClientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load embedded client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return tlsConfig, nil
}

// min helper
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
