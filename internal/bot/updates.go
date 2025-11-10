package bot

import (
	"alerting/internal/handlers"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HandleUpdates() {
	bot := InitBot()
	handlers.InitMemory()
	CheckServices(bot)
	
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
	}

	for update := range updates {
		reply := "Не знаю что сказать"
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Command() {
		case "start":
			reply = "Привет 👋 Я телеграм-бот.\n" +
				"Используй команды:\n" +
				"/add <url> — добавить ссылку\n" +
				"/list — показать сохранённые ссылки\n" +
				"/clear — очистить список"
		case "hello":
			reply = "world"
		case "add":
			args := update.Message.CommandArguments()
			if args == "" {
				reply = "Укажи URL после команды, например: /add https://example.com"
			} else {
				reply = handlers.Add(update.Message, args)
			}
		case "list":
			urls := handlers.GetURLs(update.Message.Chat.ID)
			if len(urls) == 0 {
				reply = "Ты ещё ничего не добавил 😅"
			} else {
				var sb strings.Builder
				sb.WriteString("Твои сохранённые URL:\n")
				for i, u := range urls {
					sb.WriteString(
						// нумерованный список
						strings.Join([]string{strconv.Itoa(i+1) + "️⃣ ", u, "\n"}, ""),
					)
				}
				reply = sb.String()
			}

		case "clear":
			handlers.ClearURLs(update.Message.Chat.ID)
			reply = "Все ссылки удалены 🧹"

		default:
			// если просто текст без команды
			if !update.Message.IsCommand() {
				reply = "Я понимаю команды: /add, /list, /clear"
			}
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}
