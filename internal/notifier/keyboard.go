package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

type KeyboardButton struct {
    Text string `json:"text"`
    
    // Optional fields
    RequestContact  bool `json:"request_contact,omitempty"`
    RequestLocation bool `json:"request_location,omitempty"`
    RequestPoll     *KeyboardButtonPollType `json:"request_poll,omitempty"`
    WebApp          *WebAppInfo `json:"web_app,omitempty"`
}

type KeyboardButtonPollType struct {
    Type string `json:"type,omitempty"` // "quiz" or "regular"
}

type WebAppInfo struct {
    URL string `json:"url"`
}

type ReplyKeyboardMarkup struct {
    Keyboard        [][]KeyboardButton `json:"keyboard"`
    ResizeKeyboard  bool               `json:"resize_keyboard,omitempty"`
    OneTimeKeyboard bool               `json:"one_time_keyboard,omitempty"`
    InputFieldPlaceholder string       `json:"input_field_placeholder,omitempty"`
    Selective       bool               `json:"selective,omitempty"`
}

type InlineKeyboardButton struct {
    Text string `json:"text"`
    
    // Only one of the optional fields must be used
    URL                     string `json:"url,omitempty"`
    CallbackData            string `json:"callback_data,omitempty"`
    WebApp                  *WebAppInfo `json:"web_app,omitempty"`
    LoginURL                *LoginURL `json:"login_url,omitempty"`
    SwitchInlineQuery       string `json:"switch_inline_query,omitempty"`
    SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat,omitempty"`
    CallbackGame            *CallbackGame `json:"callback_game,omitempty"`
    Pay                     bool `json:"pay,omitempty"`
}

type LoginURL struct {
    URL                string `json:"url"`
    ForwardText        string `json:"forward_text,omitempty"`
    BotUsername        string `json:"bot_username,omitempty"`
    RequestWriteAccess bool   `json:"request_write_access,omitempty"`
}

type CallbackGame struct {
    // This object represents a game. 
    // Currently empty, but exists for future extensions
}

type InlineKeyboardMarkup struct {
    InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type ReplyKeyboardRemove struct {
    RemoveKeyboard bool `json:"remove_keyboard"`
    Selective      bool `json:"selective,omitempty"`
}

type TelegramBot struct {
    token   string
    baseURL string
}

func NewTelegramBot(token string) *TelegramBot {
    return &TelegramBot{
        token:   token,
        baseURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
    }
}

// SendMessage отправка сообщения с клавиатурой
func (tb *TelegramBot) SendMessage(chatID int64, text string, replyMarkup interface{}) error {
    payload := map[string]interface{}{
        "chat_id": chatID,
        "text":    text,
    }
    
    if replyMarkup != nil {
        payload["reply_markup"] = replyMarkup
    }
    
    return tb.sendRequest("sendMessage", payload)
}

// SendReplyKeyboard отправка обычной Reply клавиатуры
func (tb *TelegramBot) SendReplyKeyboard(chatID int64, text string) error {
    keyboard := ReplyKeyboardMarkup{
        Keyboard: [][]KeyboardButton{
            {
                {Text: "📊 Статистика"},
                {Text: "⚙️ Настройки"},
            },
            {
                {Text: "📞 Контакты"},
                {Text: "ℹ️ Помощь"},
            },
        },
        ResizeKeyboard:  true,
        OneTimeKeyboard: false,
    }
    
    return tb.SendMessage(chatID, text, keyboard)
}

// SendReplyKeyboardWithContact отправка клавиатуры с кнопкой запроса контакта
func (tb *TelegramBot) SendReplyKeyboardWithContact(chatID int64, text string) error {
    keyboard := ReplyKeyboardMarkup{
        Keyboard: [][]KeyboardButton{
            {
                {Text: "📊 Статистика"},
                {Text: "⚙️ Настройки"},
            },
            {
                {Text: "📱 Поделиться контактом", RequestContact: true},
                {Text: "📍 Поделиться местоположением", RequestLocation: true},
            },
        },
        ResizeKeyboard:  true,
        OneTimeKeyboard: true,
        InputFieldPlaceholder: "Выберите действие...",
    }
    
    return tb.SendMessage(chatID, text, keyboard)
}

// SendReplyKeyboardWithPoll отправка клавиатуры с кнопкой создания опроса
func (tb *TelegramBot) SendReplyKeyboardWithPoll(chatID int64, text string) error {
    keyboard := ReplyKeyboardMarkup{
        Keyboard: [][]KeyboardButton{
            {
                {Text: "📊 Статистика"},
                {Text: "⚙️ Настройки"},
            },
            {
                {Text: "📊 Создать опрос", RequestPoll: &KeyboardButtonPollType{Type: "regular"}},
                {Text: "🧩 Создать викторину", RequestPoll: &KeyboardButtonPollType{Type: "quiz"}},
            },
        },
        ResizeKeyboard: true,
    }
    
    return tb.SendMessage(chatID, text, keyboard)
}

// SendInlineKeyboard отправка Inline клавиатуры
func (tb *TelegramBot) SendInlineKeyboard(chatID int64, text string) error {
    keyboard := InlineKeyboardMarkup{
        InlineKeyboard: [][]InlineKeyboardButton{
            {
                {Text: "✅ Да", CallbackData: "yes"},
                {Text: "❌ Нет", CallbackData: "no"},
                {Text: "⚙️ Настройки", CallbackData: "settings"},
            },
            {
                {Text: "🌐 Наш сайт", URL: "https://example.com"},
                {Text: "📱 Приложение", URL: "https://play.google.com"},
            },
            {
                {Text: "🔍 Поиск", SwitchInlineQuery: "query"},
                {Text: "💬 Чат", SwitchInlineQueryCurrentChat: "chat"},
            },
        },
    }
    
    return tb.SendMessage(chatID, text, keyboard)
}

// SendInlineKeyboardWithWebApp отправка Inline клавиатуры с WebApp
func (tb *TelegramBot) SendInlineKeyboardWithWebApp(chatID int64, text string) error {
    keyboard := InlineKeyboardMarkup{
        InlineKeyboard: [][]InlineKeyboardButton{
            {
                {Text: "📱 Открыть WebApp", WebApp: &WebAppInfo{URL: "https://your-webapp.com"}},
                {Text: "🌐 Сайт", URL: "https://example.com"},
            },
            {
                {Text: "✅ Подтвердить", CallbackData: "confirm"},
                {Text: "❌ Отмена", CallbackData: "cancel"},
            },
        },
    }
    
    return tb.SendMessage(chatID, text, keyboard)
}

// SendLoginButton отправка кнопки для авторизации
func (tb *TelegramBot) SendLoginButton(chatID int64, text string) error {
    keyboard := InlineKeyboardMarkup{
        InlineKeyboard: [][]InlineKeyboardButton{
            {
                {Text: "🔐 Войти через Telegram", 
                 LoginURL: &LoginURL{
                     URL: "https://your-site.com/auth",
                     ForwardText: "Авторизация",
                     BotUsername: "YourBot",
                     RequestWriteAccess: true,
                 }},
            },
        },
    }
    
    return tb.SendMessage(chatID, text, keyboard)
}

// RemoveKeyboard удаление клавиатуры
func (tb *TelegramBot) RemoveKeyboard(chatID int64, text string) error {
    removeKeyboard := ReplyKeyboardRemove{
        RemoveKeyboard: true,
        Selective:      false,
    }
    
    return tb.SendMessage(chatID, text, removeKeyboard)
}

// sendRequest отправка запроса к Telegram API
func (tb *TelegramBot) sendRequest(method string, payload interface{}) error {
    url := fmt.Sprintf("%s/%s", tb.baseURL, method)
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("error marshaling payload: %v", err)
    }
    
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("error sending request: %v", err)
    }
    defer resp.Body.Close()
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("error reading response: %v", err)
    }
    
    log.Printf("Response: %s", string(body))
    
    // Проверка статуса ответа
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("API error: %s", resp.Status)
    }
    
    return nil
}

func main() {
    bot := NewTelegramBot("7538507602:AAH8qYXCdK4wAn9FSoJJ4xZWm5NGtfR2Ubw")
    
    // Примеры использования
    
    // 1. Обычная Reply клавиатура
    err := bot.SendReplyKeyboard(416751006, "Выберите действие:")
    if err != nil {
        log.Printf("Error: %v", err)
    }
    
    // 2. Клавиатура с запросом контакта и местоположения
    err = bot.SendReplyKeyboardWithContact(416751006, "Поделитесь контактом или местоположением:")
    if err != nil {
        log.Printf("Error: %v", err)
    }
    
    // 3. Клавиатура с созданием опроса
    err = bot.SendReplyKeyboardWithPoll(416751006, "Создайте опрос или викторину:")
    if err != nil {
        log.Printf("Error: %v")
    }
    
    // 4. Inline клавиатура
    err = bot.SendInlineKeyboard(416751006, "Подтвердите действие:")
    if err != nil {
        log.Printf("Error: %v", err)
    }
    
    // 5. Inline клавиатура с WebApp
    err = bot.SendInlineKeyboardWithWebApp(416751006, "Откройте WebApp:")
    if err != nil {
        log.Printf("Error: %v", err)
    }
    
    // 6. Кнопка авторизации
    err = bot.SendLoginButton(416751006, "Войдите в систему:")
    if err != nil {
        log.Printf("Error: %v", err)
    }
    
    // 7. Удаление клавиатуры
    err = bot.RemoveKeyboard(416751006, "Клавиатура удалена")
    if err != nil {
        log.Printf("Error: %v", err)
    }
}