package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewEtcdClient 根据提供的配置创建etcd客户端
// 该函数会配置连接参数、认证信息和TLS设置
// 返回创建的客户端实例和可能的错误
func NewEtcdClient(cfg EtcdConfig, logger *zap.Logger) (*clientv3.Client, error) {
	log := logger.With(zap.Namespace("[etcd client]")).Sugar()
	etcConf := clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.DialTimeout,
		Username:    cfg.Username,
		Password:    cfg.Password,
		// 关键配置 2：使用 grpc.WithBlock() 确保连接建立
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithBackoffMaxDelay(3 * time.Second), // 重试间隔
		},

		// 禁用自动同步端点 - 避免获取到错误的地址（如 127.0.0.1）
		// 使用配置文件中指定的端点地址
		AutoSyncInterval: 0,

		// 关键配置 4：KeepAlive 参数
		DialKeepAliveTime:    5 * time.Second,
		DialKeepAliveTimeout: 5 * time.Second,
	}

	if cfg.TLS != nil {
		// 从第一个 endpoint 提取主机名作为 ServerName
		serverName := extractHostFromEndpoint(cfg.Endpoints[0])
		log.Info("使用 ServerName 进行 TLS 验证", zap.String("server_name", serverName))

		tlsConfig, err := createTLSConfig(cfg.TLS, serverName)
		if err != nil {
			log.Error("创建 TLS 配置失败", zap.Error(err))
			return nil, fmt.Errorf("failed to create TLS config: %w", err)
		}
		etcConf.TLS = tlsConfig
		// 使用 grpc credentials 包装 TLS 配置
		creds := credentials.NewTLS(tlsConfig)
		etcConf.DialOptions = append(etcConf.DialOptions, grpc.WithTransportCredentials(creds))
		log.Info("已启用 TLS 加密连接")
	}

	if cfg.Username != "" {
		log.Info("已启用用户认证", zap.String("username", cfg.Username))
	}

	client, err := clientv3.New(etcConf)
	if err != nil {
		log.Error("创建 etcd 客户端失败", zap.Error(err))
		return nil, err
	}

	log.Info("etcd 客户端连接成功",
		zap.Strings("endpoints", cfg.Endpoints))
	return client, nil
}

// createTLSConfig 根据 TLS 配置创建 tls.Config
func createTLSConfig(cfg *TLSConfig, serverName string) (*tls.Config, error) {
	if cfg.CAFile == "" {
		return nil, fmt.Errorf("CA file is required for TLS")
	}

	// 加载 CA 证书
	caCert, err := ioutil.ReadFile(cfg.CAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA file: %w", err)
	}

	// 创建证书池
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}

	tlsConfig := &tls.Config{
		RootCAs:    caCertPool,
		MinVersion: tls.VersionTLS12,
		ServerName: serverName, // 使用 endpoint 主机名进行证书验证
	}

	// 如果有客户端证书，加载它
	if cfg.CertFile != "" && cfg.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return tlsConfig, nil
}

// extractHostFromEndpoint 从 etcd endpoint URL 中提取主机名
// 例如: https://172.18.0.194:2379 -> 172.18.0.194
func extractHostFromEndpoint(endpoint string) string {
	// 解析 URL
	u, err := url.Parse(endpoint)
	if err != nil {
		// 如果不是标准 URL 格式，尝试直接提取 host:port
		if strings.Contains(endpoint, "://") {
			parts := strings.SplitN(endpoint, "://", 2)
			if len(parts) == 2 {
				endpoint = parts[1]
			}
		}
		// 移除端口
		if idx := strings.LastIndex(endpoint, ":"); idx != -1 {
			endpoint = endpoint[:idx]
		}
		return endpoint
	}

	// 从 Host 中移除端口
	host := u.Hostname()
	return host
}
