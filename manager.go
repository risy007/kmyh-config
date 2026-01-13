package config

import (
	"context"
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type (
	// inParams 依赖注入参数
	inParams struct {
		fx.In
		AppConfig *AppConfig
		Logger    *zap.Logger
		Client    *clientv3.Client
	}
	// ConfigManager 分布式配置管理器
	// 提供动态配置加载、监听和管理功能
	ConfigManager struct {
		client *clientv3.Client
		logger *zap.SugaredLogger
		cfg    EtcdConfig
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

// NewConfigManagerDirect 创建配置管理器（直接参数）
func NewConfigManagerDirect(client *clientv3.Client, logger *zap.Logger, appConfig *AppConfig) *ConfigManager {
	log := logger.With(zap.Namespace("[ConfigManager]")).Sugar()
	cfg := appConfig.Etcd
	return &ConfigManager{
		client: client,
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

// GetConfig 根据泛型类型自动获取配置
// 规则：将结构体名称转为小写并移除末尾的 "config" 后缀
func GetConfig[T any](m *ConfigManager, app, env string) (T, error) {
	var config T

	// 获取类型名称 - 使用零值来获取类型信息
	typeOf := reflect.TypeOf((*T)(nil)).Elem()

	// 获取基本类型（处理指针类型）
	for typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}

	groupName := strings.ToLower(typeOf.Name())

	// 如果名称以 "config" 结尾，则移除
	if strings.HasSuffix(groupName, "config") {
		groupName = groupName[:len(groupName)-6] // "config" 的长度是6
	}

	// 获取配置组
	configGroup := m.GetGroup(app, env, groupName)

	// 将配置反序列化到目标类型
	err := configGroup.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

// ConfigGroup 配置组接口
// 提供统一的配置访问方法，支持动态更新
type ConfigGroup interface {
	// Get 获取指定键的原始值
	Get(key string) interface{}
	// GetString 获取字符串类型的配置值
	GetString(key string) string
	// GetInt 获取整数类型的配置值
	GetInt(key string) int
	// GetBool 获取布尔类型的配置值
	GetBool(key string) bool
	// Unmarshal 将配置反序列化到目标对象
	Unmarshal(obj interface{}) error
	// OnChange 注册配置变更回调函数
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
	g.mu.RLock()
	defer g.mu.RUnlock()

	if val := g.viper.Get(key); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
		// 尝试将其他类型转换为字符串
		return fmt.Sprintf("%v", val)
	}
	return ""
}

// GetInt 获取整数
func (g *etcdConfigGroup) GetInt(key string) int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if val := g.viper.Get(key); val != nil {
		switch v := val.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case int32:
			return int(v)
		case float64:
			return int(v)
		case float32:
			return int(v)
		case string:
			// 尝试解析字符串为整数
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return 0
}

// GetBool 获取布尔值
func (g *etcdConfigGroup) GetBool(key string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if val := g.viper.Get(key); val != nil {
		switch v := val.(type) {
		case bool:
			return v
		case int:
			return v != 0
		case int64:
			return v != 0
		case string:
			// 尝试解析字符串为布尔值
			if b, err := strconv.ParseBool(v); err == nil {
				return b
			}
			// 处理常见的真值字符串
			lower := strings.ToLower(v)
			return lower == "true" || lower == "1" || lower == "yes" || lower == "on"
		}
	}
	return false
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
