package config

import "encoding/json"

type ChannelType string

const (
	ChannelTypeFeishu ChannelType = "feishu"
	ChannelTypeWeixin ChannelType = "weixin"
	ChannelTypeQQBot  ChannelType = "qqbot"
)

type FeishuChannelConfig struct {
	AppID             string `json:"app_id"`
	AppSecret         string `json:"app_secret"`
	CallbackType      string `json:"callback_type"`
	EncryptKey        string `json:"encrypt_key,omitempty"`
	VerificationToken string `json:"verification_token,omitempty"`
}

type WeixinChannelConfig struct {
	CorpID         string `json:"corp_id"`
	CorpSecret     string `json:"corp_secret"`
	AgentID        int64  `json:"agent_id"`
	Token          string `json:"token"`
	EncodingAESKey string `json:"encoding_aes_key,omitempty"`
	WebhookPath    string `json:"webhook_path,omitempty"`
}

type QQBotChannelConfig struct {
	AppID       string `json:"app_id"`
	AppSecret   string `json:"app_secret"`
	Token       string `json:"token"`
	WebhookPath string `json:"webhook_path,omitempty"`
	UseSandbox  bool   `json:"use_sandbox"`
}

type ChannelInitItem struct {
	TenantID    string          `json:"tenant_id"`
	ChannelID   string          `json:"channel_id"`
	ChannelType ChannelType     `json:"channel_type"`
	Status      string          `json:"status"`
	Config      json.RawMessage `json:"config"`
}

type FeishuChannelInitItem struct {
	TenantID          string `json:"tenant_id"`
	ChannelID         string `json:"channel_id"`
	AppID             string `json:"app_id"`
	AppSecret         string `json:"app_secret"`
	CallbackType      string `json:"callback_type"`
	EncryptKey        string `json:"encrypt_key,omitempty"`
	VerificationToken string `json:"verification_token,omitempty"`
	Status            string `json:"status"`
}

type WeixinChannelInitItem struct {
	TenantID    string `json:"tenant_id"`
	ChannelID   string `json:"channel_id"`
	CorpID      string `json:"corp_id"`
	CorpSecret  string `json:"corp_secret"`
	AgentID     int64  `json:"agent_id"`
	Token       string `json:"token"`
	WebhookPath string `json:"webhook_path"`
	Status      string `json:"status"`
}

type QQBotChannelInitItem struct {
	TenantID    string `json:"tenant_id"`
	ChannelID   string `json:"channel_id"`
	AppID       string `json:"app_id"`
	AppSecret   string `json:"app_secret"`
	Token       string `json:"token"`
	WebhookPath string `json:"webhook_path"`
	UseSandbox  bool   `json:"use_sandbox"`
	Status      string `json:"status"`
}
