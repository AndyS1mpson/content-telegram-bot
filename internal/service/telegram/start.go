package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const welcomeMessage = `Привет! Я собираю обои с Pinterest и публикую их в канал.

Команды:
/collect <тема> — собрать пины с Pinterest по заданной теме
/view <тема>    — показать следующий неотсмотренный пин
/publish <тема> — опубликовать все отобранные пины в канал пачкой

Пример: /collect nature 4k wallpaper`

// StartHandler обработчик команды /start
func (c *TelegramClient) StartHandler(_ context.Context, update *tgbotapi.Update) {
	if !c.validateUser(update.Message.From.ID) {
		c.sendMessage(update.Message.Chat.ID, ErrAccessDenied.Error())
		return
	}

	c.sendMessage(update.Message.Chat.ID, welcomeMessage)
}
