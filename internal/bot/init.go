package bot

import (
	"log"

	"alerting/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func InitBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(config.LoadConfig().TelegramToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account: %s", bot)

	return bot
}
