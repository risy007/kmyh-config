package config

// TaskConfig 模块级任务配置
type TaskConfig struct {
	// 模块名称
	Module string `mapstructure:"module"`
	// 任务列表
	Jobs []JobConfig `mapstructure:"jobs"`
}

// JobConfig 单个任务配置
type JobConfig struct {
	// 任务名称（必须唯一）
	Name string `mapstructure:"name"`
	// Cron 表达式
	Spec string `mapstructure:"spec"`
	// 任务描述
	Description string `mapstructure:"desc"`
	// 是否启用
	Enabled bool `mapstructure:"enabled"`
	// 额外参数
	Params map[string]interface{} `mapstructure:"params"`
}

// CronJobConfig 兼容旧版配置
type CronJobConfig struct {
	Name        string `mapstructure:"name"`
	Spec        string `mapstructure:"spec"`
	Description string `mapstructure:"desc"`
	SendWXMQ    bool   `mapstructure:"send_wxmq"`
	Enabled     bool   `mapstructure:"enabled"`
}

// TaskConfigList 兼容旧版配置
type TaskConfigList []CronJobConfig
