package config

import (
	"os"
	"log"
	"path/filepath"
	"alerting/internal/models"
	"gopkg.in/yaml.v3"
)

var Cfg *models.Config

func Init() {
    LoadConfig("internal/config/config.yaml")
}

func LoadConfig(path string) {
    absPath, err := filepath.Abs(path)
    if err != nil {
        log.Fatalf("Ошибка при обработке пути: %v", err)
    }

	data, err := os.ReadFile(absPath)
	if err != nil {
		log.Fatalf("Ошибка чтения файла конфигурации: %v", err)
	}

	var cfg models.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Ошибка парсинга YAML: %v", err)
	}

	Cfg = &cfg
}
