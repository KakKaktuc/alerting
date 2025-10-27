package notifier

import (
	"alerting/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

type Chat struct {
	ID int64 `json:"id"`
}

func GetChatID() int64 {
	token := config.LoadConfig().TelegramToken
	url := fmt.Sprintf("https://api.telegram.org/bot%v/getUpdates", token)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Ошибка при запросе:", err)
		return 0
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка при чтении тела:", err)
		return 0
	}

	var data UpdateResponse
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("Ошибка при парсинге JSON:", err)
		return 0
	}

	// Проверяем, что result не пустой
	if len(data.Result) == 0 {
		log.Println("Нет обновлений.")
		return 0
	}

	// Берём chat.id из последнего апдейта
	chatID := data.Result[len(data.Result)-1].Message.Chat.ID
	return chatID
}
