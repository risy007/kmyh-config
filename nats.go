package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ---------------------------------------------------------------------------
// 连接工厂
// ---------------------------------------------------------------------------

// NewNatsConnection 根据配置创建 NATS 连接
// appName 用作客户端名称（显示在 NATS 服务端连接列表，auth callout 识别客户端）
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

	if needNatsTLS(cfg) {
		tlsConfig, err := createTLSConfig(cfg.TLS, "")
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

func needNatsTLS(cfg NatsConfig) bool {
	return cfg.TLS != nil || strings.HasPrefix(cfg.Address, "tls://")
}

// ---------------------------------------------------------------------------
// TLS 配置（etcd + NATS 共用）
// ---------------------------------------------------------------------------

// createTLSConfig 根据 TLS 配置创建 tls.Config
// 如果 cfg 为 nil 或字段为空，使用 embedded 证书
// serverName 留空时跳过 ServerName 验证（NATS 场景）
func createTLSConfig(cfg *TLSConfig, serverName string) (*tls.Config, error) {
	caCert := embeddedCA
	if cfg != nil && cfg.CAFile != "" {
		data, err := os.ReadFile(cfg.CAFile)
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
	if serverName != "" {
		tlsConfig.ServerName = serverName
	}

	if cfg != nil && cfg.CertFile != "" && cfg.KeyFile != "" {
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

// ---------------------------------------------------------------------------
// 托管服务
// ---------------------------------------------------------------------------

type IMicroService interface {
	RegisterMicro(nc *nats.Conn) error
	StopMicro()
}

type Service struct {
	log     *zap.SugaredLogger
	cfg     *NatsConfig
	appName string
	logger  *zap.Logger

	mu        sync.RWMutex
	nc        *nats.Conn
	stop      chan struct{}
	microSvcs []IMicroService
}

func (s *Service) connect() error {
	if s.nc != nil {
		return nil
	}

	nc, err := NewNatsConnection(*s.cfg, s.appName, s.logger)
	if err != nil {
		s.log.Error("NATS 连接失败", zap.Error(err))
		return err
	}
	s.nc = nc
	s.log.Info("NATS 初始化连接成功", zap.String("address", s.cfg.Address))
	return nil
}

func (s *Service) retryConnect() {
	if s.stop == nil {
		return
	}
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stop:
			s.log.Info("NATS 后台重试已停止")
			return
		case <-ticker.C:
			s.mu.Lock()
			if s.nc != nil {
				s.mu.Unlock()
				return
			}
			nc, err := NewNatsConnection(*s.cfg, s.appName, s.logger)
			if err != nil {
				s.mu.Unlock()
				s.log.Warn("NATS 后台重连失败，继续重试...", zap.Error(err))
				continue
			}
			s.nc = nc
			svcs := s.microSvcs
			s.mu.Unlock()
			s.log.Info("NATS 后台重连成功", zap.String("address", s.cfg.Address))
			s.registerAllMicro(nc, svcs)
			return
		}
	}
}

func (s *Service) Reload(newCfg *NatsConfig) error {
	s.log.Info("正在重新加载 NATS 连接...", zap.String("address", newCfg.Address))

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.nc != nil {
		s.nc.Close()
		s.nc = nil
	}

	s.cfg = newCfg
	if err := s.connect(); err != nil {
		go s.retryConnect()
		return err
	}

	s.registerAllMicro(s.nc, s.microSvcs)
	return nil
}

func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-s.stop:
	default:
		close(s.stop)
	}

	for _, svc := range s.microSvcs {
		svc.StopMicro()
	}

	if s.nc != nil {
		s.nc.Close()
		s.nc = nil
		s.log.Info("NATS 连接已关闭")
	}
	return nil
}

func (s *Service) GetConnection() *nats.Conn {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.nc
}

func (s *Service) registerAllMicro(nc *nats.Conn, svcs []IMicroService) {
	for _, svc := range svcs {
		if err := svc.RegisterMicro(nc); err != nil {
			s.log.Error("注册 Micro 服务失败", zap.Error(err))
		}
	}
}

// ---------------------------------------------------------------------------
// FX 模块
// ---------------------------------------------------------------------------

type natsInParams struct {
	fx.In
	Logger    *zap.Logger
	CfgMgr    *ConfigManager
	AppConfig *AppConfig
}

type natsOutResult struct {
	fx.Out
	Service *Service
	NC      *nats.Conn
}

func NewNatsService(in natsInParams) (natsOutResult, error) {
	log := in.Logger.With(zap.Namespace("[NatsMQ-Service]")).Sugar()

	s := &Service{
		log:     log,
		cfg:     &NatsConfig{},
		appName: in.AppConfig.AppName,
		logger:  in.Logger,
		stop:    make(chan struct{}),
	}

	CfgGroup, err := in.CfgMgr.GetGroup(in.AppConfig.AppName, in.AppConfig.Env, "nats")
	if err != nil {
		log.Error("获取 NATS 配置失败，模块启动中止", zap.Error(err))
		return natsOutResult{}, fmt.Errorf("failed to get nats config: %w", err)
	}

	CfgGroup.OnChange(func() {
		log.Info("NATS 配置已更新，需要重新连接")
		var newConf NatsConfig
		if err := CfgGroup.Unmarshal(&newConf); err != nil {
			log.Warn("解析 NATS 配置失败", zap.Error(err))
			return
		}
		if err := s.Reload(&newConf); err != nil {
			log.Error("NATS 重连失败", zap.Error(err))
		}
	})

	if err := CfgGroup.Unmarshal(s.cfg); err != nil {
		log.Error("解析 NATS 配置失败", zap.Error(err))
		return natsOutResult{}, fmt.Errorf("failed to unmarshal nats config: %w", err)
	}

	s.mu.Lock()
	if err := s.connect(); err != nil {
		s.mu.Unlock()
		log.Warn("启动 NATS 连接失败，将在后台重试", zap.Error(err))
		go s.retryConnect()
	} else {
		s.mu.Unlock()
	}

	return natsOutResult{
		Service: s,
		NC:      s.GetConnection(),
	}, nil
}

type microServicesParams struct {
	fx.In
	Service *Service
	Micros  []IMicroService `group:"micro_services"`
}

func RegisterMicroServices(in microServicesParams) {
	s := in.Service
	s.mu.Lock()
	s.microSvcs = in.Micros
	nc := s.nc
	s.mu.Unlock()

	if nc != nil {
		s.registerAllMicro(nc, in.Micros)
	}
}

func NewNatsModule() fx.Option {
	return fx.Module("natsmq",
		fx.Provide(NewNatsService),
		fx.Invoke(RegisterMicroServices),
	)
}
