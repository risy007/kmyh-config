package config

import (
	"crypto/tls"
	"google.golang.org/grpc"
	"time"

	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
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

		// 关键配置 3：启用自动重连
		AutoSyncInterval: 30 * time.Second, // 自动同步端点列表

		// 关键配置 4：KeepAlive 参数
		DialKeepAliveTime:    5 * time.Second,
		DialKeepAliveTimeout: 5 * time.Second,
	}

	if cfg.TLS != nil {
		tlsConfig, err := tls.LoadX509KeyPair(cfg.TLS.CertFile, cfg.TLS.KeyFile)
		if err != nil {
			return nil, err
		}
		etcConf.TLS = &tls.Config{
			Certificates: []tls.Certificate{tlsConfig},
		}
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
