package main

import (
	"time"

	"alerting/internal/notifier"
	"alerting/internal/checker"
)


for_, url := range urls {
	go func(u string) {
		for {
			if err :+ checker.CheckURL(u); err != nil {
				notifier.SendTelegram(fmt.Sprintf("%s down: %v", u, err))
			}
			time.Sleep(interval)
		}
	}(url)
}