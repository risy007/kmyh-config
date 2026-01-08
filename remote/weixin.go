package config

type (
	WorkwxWebHookConfig struct {
		Key       string `mapstructure:"key"`
		Subscribe string `mapstructure:"subject"`
	}
	WorkwxAppConfig struct {
		Address        string `mapstructure:"address"`
		CorpSecret     string `mapstructure:"corp_secret"`
		AgentID        int64  `mapstructure:"agent_id"`
		Token          string `mapstructure:"token"`
		EncodingAESKey string `mapstructure:"encoding_AESKey"`
		TxSubscribe    string `mapstructure:"tx_subject"`
		RxSubscribe    string `mapstructure:"rx_subject"`
	}
	WeixinConfig struct {
		Enabled           bool                `mapstructure:"enabled"`
		CorpID            string              `mapstructure:"corp_id"`
		WebHook           WorkwxWebHookConfig `mapstructure:"WebHook"`
		App               WorkwxAppConfig     `mapstructure:"App"`
		QYAPIHostOverride string              `mapstructure:"QYAPIHostOverride"`
		TLSKeyLogFile     string              `mapstructure:"TLSKeyLogFile"`
	}
)
