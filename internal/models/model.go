package models

type Config struct {
	APIKey string `yaml:"telegram_token"`
	ChatID string `yaml:"telegram_chat_id"`

	URLs []string `yaml:"urls"`
}

