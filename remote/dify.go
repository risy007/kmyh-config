package config

type (
	DifyConfig struct {
		BaseURL       string `mapstructure:"base_url"`
		APIKey        string `mapstructure:"api_key"`
		CachePeriod   string `mapstructure:"cache_period"`
		DefaultPrompt string `mapstructure:"default_prompt"`
		BotType       string `mapstructure:"bot_type"`
		WorkflowID    string `mapstructure:"workflow_id"`
	}
)
