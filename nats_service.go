package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type IMicroService interface {
	RegisterMicro(nc *nats.Conn) error
	StopMicro()
}

type (
	natsInParams struct {
		fx.In
		Logger    *zap.Logger
		CfgMgr    *ConfigManager
		AppConfig *AppConfig
	}

	natsOutResult struct {
		fx.Out
		Service *Service
		NC      *nats.Conn
	}

	Service struct {
		log     *zap.SugaredLogger
		cfg     *NatsConfig
		appName string
		logger  *zap.Logger

		mu        sync.RWMutex
		nc        *nats.Conn
		stop      chan struct{}
		microSvcs []IMicroService
	}
)

func (s *Service) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.connect()
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

func (s *Service) Publish(subject string, data []byte) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.nc == nil {
		return nats.ErrConnectionClosed
	}
	return s.nc.Publish(subject, data)
}

func (s *Service) Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.nc == nil {
		return nil, nats.ErrConnectionClosed
	}
	return s.nc.Subscribe(subject, handler)
}

func (s *Service) QueueSubscribe(subject, queue string, handler nats.MsgHandler) (*nats.Subscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.nc == nil {
		return nil, nats.ErrConnectionClosed
	}
	return s.nc.QueueSubscribe(subject, queue, handler)
}

func (s *Service) Request(subject string, data []byte, timeout int) (*nats.Msg, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.nc == nil {
		return nil, nats.ErrConnectionClosed
	}
	timeoutDuration := time.Duration(timeout) * time.Second
	if timeoutDuration <= 0 {
		timeoutDuration = 5 * time.Second
	}
	return s.nc.Request(subject, data, timeoutDuration)
}

func (s *Service) registerAllMicro(nc *nats.Conn, svcs []IMicroService) {
	for _, svc := range svcs {
		if err := svc.RegisterMicro(nc); err != nil {
			s.log.Error("注册 Micro 服务失败", zap.Error(err))
		}
	}
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

	if err := s.Start(); err != nil {
		log.Warn("启动 NATS 连接失败，将在后台重试", zap.Error(err))
		go s.retryConnect()
	}

	return natsOutResult{
		Service: s,
		NC:      s.GetConnection(),
	}, nil
}
