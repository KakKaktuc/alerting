package handlers

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	userMemory = make(map[int64][]string)
	mu         sync.Mutex
)

func saveMemory() error {
	file, err := os.Create("memory.json")
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(userMemory)
}

func loadMemory() error {
	file, err := os.Open("memory.json")
	log.Println(file)
	if os.IsNotExist(err) {
		userMemory = make(map[int64][]string)
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(&userMemory)
}

func InitMemory() {
	if err := loadMemory(); err != nil {
		panic("Не удалось загрузить память: " + err.Error())
	}
}

func isValidURL(raw string) bool {
	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}
	if parsed.Scheme != "https" {
		return false
	}
	if parsed.Host == "" || strings.Contains(parsed.Host, " ") {
		return false
	}
	return true
}

func Add(msg *tgbotapi.Message, url string) (reply string) {
	if !isValidURL(url) {
		return "Неверный формат ссылки. Используй формат: https://example.com"
	}
	mu.Lock()
	defer mu.Unlock()

	chatID := msg.Chat.ID
	userMemory[chatID] = append(userMemory[chatID], url)

	_ = saveMemory()

	return "Ссылка добавлена"
}

func GetURLs(ChatID int64) []string {
	mu.Lock()
	defer mu.Unlock()

	return append([]string(nil), userMemory[ChatID]...)
}

func ClearURLs(chatID int64) {
	mu.Lock()
	defer mu.Unlock()

	delete(userMemory, chatID)
	_ = saveMemory()
}
