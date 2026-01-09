package config

import (
	"context"
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"github.com/risy007/kmyh-config/local"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sync"
)

type (
	inParams struct {
		fx.In
		AppConfig *AppConfig
		Logger    *zap.Logger
		Client    *clientv3.Client
	}
	// ConfigManager 配置管理器
	ConfigManager struct {
		client *clientv3.Client
		logger *zap.SugaredLogger
		cfg    config.EtcdConfig
		groups map[string]ConfigGroup
		mu     sync.RWMutex
	}
)

// NewConfigManager 创建配置管理器
func NewConfigManager(in inParams) *ConfigManager {
	log := in.Logger.With(zap.Namespace("[ConfigManager]")).Sugar()
	cfg := in.AppConfig.Etcd
	return &ConfigManager{
		client: in.Client,
		logger: log,
		groups: make(map[string]ConfigGroup),
		cfg:    cfg,
	}
}

// GetGroup 获取配置组（不存在则创建）
func (m *ConfigManager) GetGroup(app, env, group string) ConfigGroup {
	key := fmt.Sprintf("%s/%s/%s/%s/content.yaml", m.cfg.Prefix, app, env, group)

	m.mu.RLock()
	if g, exists := m.groups[key]; exists {
		m.mu.RUnlock()
		return g
	}
	m.mu.RUnlock()

	// 创建新配置组
	v := viper.New()
	v.SetConfigType("yaml")

	// 添加 etcd 远程提供者
	if err := v.AddRemoteProvider("etcd3",
		m.client.Endpoints()[0], key); err != nil {
		m.logger.Fatal("添加远程提供者失败", zap.Error(err))
	}

	// 初始读取
	if err := v.ReadRemoteConfig(); err != nil {
		m.logger.Warn("读取远程配置失败，使用空配置",
			zap.String("key", key), zap.Error(err))
	}

	g := &etcdConfigGroup{
		viper:    v,
		logger:   m.logger.With(zap.String("group", key)),
		groupKey: key,
		watchers: []func(){},
	}

	// 注册到管理器
	m.mu.Lock()
	m.groups[key] = g
	m.mu.Unlock()

	// 启动该配置组的动态监听
	go m.watchGroup(g)

	return g
}

// watchGroup 监听配置组变更
func (m *ConfigManager) watchGroup(g *etcdConfigGroup) {
	watchKey := g.groupKey + "/"

	watchChan := m.client.Watch(context.Background(), watchKey,
		clientv3.WithPrefix())

	g.logger.Info("开始监听配置变更", zap.String("watch_key", watchKey))

	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			g.logger.Info("配置变更事件",
				zap.String("key", string(event.Kv.Key)),
				zap.String("value", string(event.Kv.Value)))

			// 重新读取配置
			if err := g.viper.ReadRemoteConfig(); err != nil {
				g.logger.Error("重新读取配置失败", zap.Error(err))
				continue
			}

			// 通知监听者
			g.notifyWatchers()
		}
	}
}

// StartWatching 启动所有配置组的监听（fx.Invoke 调用）
func (m *ConfigManager) StartWatching() {
	m.logger.Info("配置管理器启动动态监听")
}

// Stop 停止所有监听（fx 生命周期）
func (m *ConfigManager) Stop(ctx context.Context) error {
	m.logger.Info("配置管理器停止")
	return m.client.Close()
}

// ConfigGroup 配置组接口
type ConfigGroup interface {
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	Unmarshal(obj interface{}) error
	OnChange(fn func())
}

// etcdConfigGroup 基于 etcd 的配置组实现
type etcdConfigGroup struct {
	viper    *viper.Viper
	logger   *zap.SugaredLogger
	groupKey string // 例如: /configs/myapp/prod/database
	watchers []func()
	mu       sync.RWMutex
}

// Get 获取原始值
func (g *etcdConfigGroup) Get(key string) interface{} {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.viper.Get(key)
}

// GetString 获取字符串
func (g *etcdConfigGroup) GetString(key string) string {
	return g.Get(key).(string)
}

// GetInt 获取整数
func (g *etcdConfigGroup) GetInt(key string) int {
	return g.Get(key).(int)
}

// GetBool 获取布尔值
func (g *etcdConfigGroup) GetBool(key string) bool {
	return g.Get(key).(bool)
}

// Unmarshal 反序列化到结构体
func (g *etcdConfigGroup) Unmarshal(obj interface{}) error {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.viper.Unmarshal(obj, func(config *mapstructure.DecoderConfig) {
		config.TagName = "mapstructure" // 使用 mapstructure tag
	})
}

// OnChange 注册配置变更回调
func (g *etcdConfigGroup) OnChange(fn func()) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.watchers = append(g.watchers, fn)
}

// notifyWatchers 通知所有监听者
func (g *etcdConfigGroup) notifyWatchers() {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, fn := range g.watchers {
		go fn() // 异步执行避免阻塞
	}
}
