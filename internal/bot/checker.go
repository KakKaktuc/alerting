package bot

import (
	"log"
	"net/http"
	"time"

	"alerting/internal/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// структура для хранения состояния сайта
type serviceStatus struct {
	isDown      bool
	lastNotified time.Time
}

func CheckServices(bot *tgbotapi.BotAPI) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// карта для статусов: url -> состояние
	statusMap := make(map[string]*serviceStatus)

	go func() {
		for {
			time.Sleep(60 * time.Second)

			urlsByUser := handlers.GetAllURLs() // 🔁 всегда актуальные ссылки

			for chatID, urls := range urlsByUser {
				for _, u := range urls {
					req, _ := http.NewRequest("GET", u, nil)
					req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ServiceChecker/1.0; +https://example.com)")

					resp, err := client.Do(req)

					isDown := err != nil || resp.StatusCode >= 400
					if resp != nil {
						resp.Body.Close()
					}

					st, exists := statusMap[u]
					if !exists {
						st = &serviceStatus{}
						statusMap[u] = st
					}

					// 🚨 Если сервис упал
					if isDown && !st.isDown {
						msg := tgbotapi.NewMessage(chatID, "🚨 Сервис недоступен: "+u)
						bot.Send(msg)
						log.Printf("[DOWN] %s недоступен (код: %v)", u, resp)
						st.isDown = true
						st.lastNotified = time.Now()
						continue
					}

					// ⚠️ Если сервис лежит дольше 2 минут — напомнить
					if isDown && time.Since(st.lastNotified) > 2*time.Minute {
						msg := tgbotapi.NewMessage(chatID, "⏰ Сервис всё ещё недоступен: "+u)
						bot.Send(msg)
						st.lastNotified = time.Now()
						continue
					}

					// ✅ Если восстановился
					if !isDown && st.isDown {
						msg := tgbotapi.NewMessage(chatID, "✅ Сервис снова доступен: "+u)
						bot.Send(msg)
						log.Printf("[UP] %s восстановился", u)
						st.isDown = false
						st.lastNotified = time.Now()
					}
				}
			}
		}
	}()
}
