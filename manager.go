package config

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
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
// 如果远程配置获取失败，返回错误
func (m *ConfigManager) GetGroup(app, env, group string) (ConfigGroup, error) {
	groupKey := fmt.Sprintf("%s/%s/%s/%s", m.cfg.Prefix, app, env, group)
	configKey := groupKey + "/content.yaml"

	m.mu.RLock()
	if g, exists := m.groups[groupKey]; exists {
		m.mu.RUnlock()
		return g, nil
	}
	m.mu.RUnlock()

	// 直接使用 etcd 客户端读取配置（支持 TLS 和认证）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := m.client.Get(ctx, configKey)
	if err != nil {
		m.logger.Error("从 etcd 读取配置失败",
			zap.String("key", configKey), zap.Error(err))
		return nil, fmt.Errorf("failed to read config from etcd for %s: %w", configKey, err)
	}

	if len(resp.Kvs) == 0 {
		m.logger.Error("配置键不存在",
			zap.String("key", configKey))
		return nil, fmt.Errorf("config key not found: %s", configKey)
	}

	// 创建 viper 实例并加载配置
	v := viper.New()
	v.SetConfigType("yaml")

	// 将 etcd 中的值加载到 viper
	if err := v.ReadConfig(bytes.NewReader(resp.Kvs[0].Value)); err != nil {
		m.logger.Error("解析配置失败",
			zap.String("key", configKey), zap.Error(err))
		return nil, fmt.Errorf("failed to parse config for %s: %w", configKey, err)
	}

	m.logger.Info("成功加载配置",
		zap.String("key", configKey),
		zap.Int("size", len(resp.Kvs[0].Value)))

	g := &etcdConfigGroup{
		viper:    v,
		client:   m.client,
		logger:   m.logger.With(zap.String("group", groupKey)),
		groupKey: groupKey,
		watchers: []func(){},
	}

	// 注册到管理器
	m.mu.Lock()
	m.groups[groupKey] = g
	m.mu.Unlock()

	// 启动该配置组的动态监听
	go m.watchGroup(g)

	return g, nil
}

// watchGroup 监听配置组变更
func (m *ConfigManager) watchGroup(g *etcdConfigGroup) {
	watchKey := g.groupKey + "/"
	configKey := g.groupKey + "/content.yaml"

	watchChan := m.client.Watch(context.Background(), watchKey,
		clientv3.WithPrefix())

	g.logger.Info("开始监听配置变更", zap.String("watch_key", watchKey))

	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			g.logger.Info("配置变更事件",
				zap.String("key", string(event.Kv.Key)),
				zap.String("value", string(event.Kv.Value)))

			// 重新读取配置 - 使用 etcd 客户端直接读取
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			resp, err := g.client.Get(ctx, configKey)
			cancel()

			if err != nil {
				g.logger.Error("重新读取配置失败", zap.Error(err))
				continue
			}

			if len(resp.Kvs) == 0 {
				g.logger.Error("配置键不存在", zap.String("key", configKey))
				continue
			}

			// 重新加载配置到 viper
			g.mu.Lock()
			if err := g.viper.ReadConfig(bytes.NewReader(resp.Kvs[0].Value)); err != nil {
				g.mu.Unlock()
				g.logger.Error("解析配置失败", zap.Error(err))
				continue
			}
			g.mu.Unlock()

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
// 如果远程配置获取失败，返回错误
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
	configGroup, err := m.GetGroup(app, env, groupName)
	if err != nil {
		return config, fmt.Errorf("failed to get config group for %s/%s/%s: %w", app, env, groupName, err)
	}

	// 将配置反序列化到目标类型
	if err := configGroup.Unmarshal(&config); err != nil {
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
	client   *clientv3.Client // etcd 客户端，用于重新加载配置
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
