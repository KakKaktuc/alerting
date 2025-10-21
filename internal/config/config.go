package config

import (
	"os"
	"fmt"
	"gopkg.in/yaml.v3"
)

type TelegramConfig struct {
	APIKey string `yaml:"telegram_token"`
	ChatID string `yaml:"telegram_chat_id"`
}

var AppConfig *TelegramConfig

func Init() error {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config TelegramConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}
	AppConfig = &config
	return nil
}

func Get() *TelegramConfig {
	return AppConfig
}