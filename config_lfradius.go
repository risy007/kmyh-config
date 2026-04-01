package config

type LFRadiusConfig struct {
	API          string          `mapstructure:"api"`
	BaseURL      string          `mapstructure:"base_url"`
	TimeFormat   string          `mapstructure:"time_format"`
	Remote       bool            `mapstructure:"remote"`
	BackupDir    string          `mapstructure:"backup_dir"`
	ProcessedDir string          `mapstructure:"processed_dir"`
	Tasks        []CronJobConfig `mapstructure:"tasks"`
}
