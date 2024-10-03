package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const welcomeMessage = `
	Привет! Это бот для ведения каналов с получением контента и публикацией контента в каналы.
`

// StartHandler обработчик команды /start
func (c *TelegramClient) StartHandler(_ context.Context, update *tgbotapi.Update) {
	if !c.validateUser(update.Message.From.ID) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, ErrAccessDenied.Error())
		c.bot.Send(msg)

		return
	}

	// кнопки с командами
	parsePinterestButton := tgbotapi.NewKeyboardButton(string(CommandParsePinterest))

	// создаем клавиатуру
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(parsePinterestButton),
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeMessage)
	msg.ReplyMarkup = keyboard
	c.bot.Send(msg)
}
