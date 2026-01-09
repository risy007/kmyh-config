package config

import (
	"github.com/risy007/kmyh-config/local"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Module("config",
		// 提供配置结构（从 viper 读取）
		fx.Provide(
			config.NewAppConfig,  // 加载主配置文件
			config.NewZapLogger,  // 创建日志记录器
			config.NewEtcdClient, // 创建 etcd 客户端
			NewConfigManager,     // 创建配置管理器
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
