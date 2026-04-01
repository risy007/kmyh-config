package config

import "encoding/json"

type AgentType string

const (
	AgentTypeOpenclaw AgentType = "openclaw"
	AgentTypeGoclaw   AgentType = "goclaw"
	AgentTypeDify     AgentType = "dify"
)

type AgentConfig struct {
	Type     AgentType          `json:"agent_type"`
	Name     string             `json:"name,omitempty"`
	Openclaw *OpenclawBotConfig `json:"openclaw,omitempty"`
	Goclaw   *GoclawBotConfig   `json:"goclaw,omitempty"`
	Dify     *DifyBotConfig     `json:"dify,omitempty"`
}

func (c *AgentConfig) MarshalJSON() ([]byte, error) {
	type Alias AgentConfig
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

func (c *AgentConfig) UnmarshalJSON(data []byte) error {
	type Alias AgentConfig
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	return nil
}

type CustomerBindInfo struct {
	TenantID      string `json:"tenant_id"`
	UnifiedUserID string `json:"unified_user_id"`
	AgentID       string `json:"agent_id"`
	BaseURL       string `json:"base_url"`
	APIKey        string `json:"api_key"`
}

type GetCustomerBindInfoRequest struct {
	Channel       string `json:"channel"`
	ChannelUserID string `json:"channel_user_id"`
}

type GetCustomerBindInfoResponse struct {
	Found bool              `json:"found"`
	Info  *CustomerBindInfo `json:"info,omitempty"`
	Error string            `json:"error,omitempty"`
}

type NATSResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type NATSSuccessResponse struct {
	Data interface{} `json:"data"`
}

type NATSErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	NATSCodeSuccess            = 0
	NATSCodeNotFound           = 1001
	NATSCodeUnauthorized       = 1002
	NATSCodeChannelDisabled    = 1003
	NATSCodeChannelUnsupported = 1004
	NATSCodeInternalError      = 2001
)
