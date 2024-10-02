package telegram

import (
	"github.com/pkg/errors"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// NewTgClient конструктор для клиента к API Telegram
func NewTgClient(token string) (*tgBotApi.BotAPI, error) {
	bot, err := tgBotApi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create telegram bot")
	}

	bot.Debug = false

	return bot, nil
}
