module taskbot

go 1.22.2

require github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible

require github.com/technoweenie/multipartstreamer v1.0.1 // indirect

replace github.com/go-telegram-bot-api/telegram-bot-api => ./local/telegram-bot-api
