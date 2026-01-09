package pkg

import (
	"crypto/tls"
	"github.com/risy007/kmyh-config"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"time"

	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type etcdClientParams struct {
	fx.In
	AppConfig *config.AppConfig
	Logger    *zap.Logger
}

// NewEtcdClient 创建 etcd 客户端（fx 提供者）
func NewEtcdClient(in etcdClientParams) (*clientv3.Client, error) {
	cfg := in.AppConfig.Etcd
	log := in.Logger.With(zap.Namespace("[etcd client]")).Sugar()
	etcConf := clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.DialTimeout,
		Username:    cfg.Username,
		Password:    cfg.Password,
		// 关键配置 2：使用 grpc.WithBlock() 确保连接建立[^50^]
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
