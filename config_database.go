package config

import (
	"fmt"
	"net/url"
	"strings"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Engine      string `mapstructure:"engine" description:"数据库引擎"`
	Name        string `mapstructure:"name" description:"数据库名称"`
	Host        string `mapstructure:"host" description:"服务器地址"`
	Port        int    `mapstructure:"port" description:"端口号"`
	Username    string `mapstructure:"username" description:"用户名"`
	Password    string `mapstructure:"password" description:"密码"`
	TablePrefix string `mapstructure:"table_prefix" description:"表前缀"`
	Parameters  string `mapstructure:"parameters" description:"连接参数"`

	MaxLifetime  int `mapstructure:"max_lifetime" description:"连接最大生命周期(秒)"`
	MaxOpenConns int `mapstructure:"max_open_conns" description:"最大打开连接数"`
	MaxIdleConns int `mapstructure:"max_idle_conns" description:"最大空闲连接数"`
}

// Validate 验证数据库配置
func (cfg *DatabaseConfig) Validate() error {
	if cfg.Engine == "" {
		return fmt.Errorf("database engine is required")
	}
	if cfg.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if cfg.Port <= 0 {
		return fmt.Errorf("database port must be greater than 0")
	}
	if cfg.Name == "" {
		return fmt.Errorf("database name is required")
	}
	return nil
}

func (cfg *DatabaseConfig) Dsn() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	// 添加额外参数
	params := cfg.parseParameters()
	if len(params) > 0 {
		dsn += "?" + params.Encode()
	}

	return dsn
}

// parseParameters 解析参数字符串为 url.Values
// 支持格式: "param1=value1&param2=value2" 或 "param1=value1,param2=value2"
func (cfg *DatabaseConfig) parseParameters() url.Values {
	values := url.Values{}

	if cfg.Parameters == "" {
		return values
	}

	// 支持 & 或 , 分隔符
	separator := "&"
	if !strings.Contains(cfg.Parameters, "&") && strings.Contains(cfg.Parameters, ",") {
		separator = ","
	}

	pairs := strings.Split(cfg.Parameters, separator)
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			values.Add(kv[0], kv[1])
		}
	}

	return values
}
