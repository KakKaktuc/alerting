package config

import (
	"alerting/internal/models"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() *models.Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл")
	}
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN не найден в .env")
	}

	return &models.Config{
		TelegramToken: token,
	}
}