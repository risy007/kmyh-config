package config

import (
	"fmt"
	"net/url"
	"strings"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Engine      string `mapstructure:"engine"`
	Name        string `mapstructure:"name"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	TablePrefix string `mapstructure:"table_prefix"`
	Parameters  string `mapstructure:"parameters"`

	MaxLifetime  int `mapstructure:"max_lifetime"`
	MaxOpenConns int `mapstructure:"max_open_conns"`
	MaxIdleConns int `mapstructure:"max_idle_conns"`
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
