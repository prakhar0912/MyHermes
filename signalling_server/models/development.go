package models

type Deployment struct {
	ListenAt          string `toml:"listen_at"`
	LeakServerAddress string `toml:"leak_server_address"`
	BotUID            string `toml:"bot_uid"`
	BotSecret         string `toml:"bot_secret"`
	SelfWebhook       string `toml:"self_webhook"`
}
