package config

type CasbinConfig struct {
	Enable             bool     `mapstructure:"enable" description:"是否启用权限控制"`
	Debug              bool     `mapstructure:"debug" description:"调试模式"`
	Model              string   `mapstructure:"model" description:"模型配置文件路径"`
	AutoLoad           bool     `mapstructure:"auto_load" description:"自动加载策略"`
	AutoLoadInternal   int      `mapstructure:"auto_load_internal" description:"自动加载间隔(秒)"`
	IgnorePathPrefixes []string `mapstructure:"ignore_path_prefixes" description:"忽略的路径前缀"`
}
