package config

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewConfigModule 创建FX模块（保留原有功能）
func NewConfigModule() fx.Option {
	return fx.Options(
		fx.Module("config",
			// 提供配置结构（从 viper 读取）
			fx.Provide(
				NewAppConfig,               // 加载主配置文件
				newZapLoggerFromAppConfig,  // 从AppConfig创建日志记录器
				newEtcdClientFromAppConfig, // 从AppConfig创建etcd客户端
				NewConfigManager,           // 创建配置管理器
			),
			// 生命周期管理
			fx.Invoke(
				func(manager *ConfigManager) {
					manager.StartWatching()
				},
			),
			fx.Decorate(
				// 自动注册停止钩子
				func(lifecycle fx.Lifecycle, manager *ConfigManager) *ConfigManager {
					lifecycle.Append(fx.Hook{
						OnStop: manager.Stop,
					})
					return manager
				},
			),
		))
}

// 内部辅助函数
func newZapLoggerFromAppConfig(cfg *AppConfig) *zap.Logger {
	return NewZapLogger(cfg.Logger)
}

func newEtcdClientFromAppConfig(cfg *AppConfig, logger *zap.Logger) (*clientv3.Client, error) {
	return NewEtcdClient(cfg.Etcd, logger)
}
