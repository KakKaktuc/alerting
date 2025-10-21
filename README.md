# alerting
В папку internal/config нужно добавить файл `config.yaml` в котором нужно указать ссылки 
примерный формат файла:
```
urls:
  - https://google.com
  - https://youtube.com
interval: 60s
telegram_token: "Token"
telegram_chat_id: "<chatID>"
```
