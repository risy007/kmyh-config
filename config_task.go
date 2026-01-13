package config

type CronJobConfig struct {
	Name        string `mapstructure:"name"`
	Spec        string `mapstructure:"spec"`
	Description string `mapstructure:"desc"`
	SendWXMQ    bool   `mapstructure:"send_wxmq"`
	Enabled     bool   `mapstructure:"enabled"`
}

type TaskConfig []CronJobConfig
