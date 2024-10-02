package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ParsePinsHandler обработчик команды запуска парсинга картинок
func (c *TelegramClient) ParsePinsHandler(ctx context.Context, update *tgbotapi.Update) {
	if !c.validateUser(update.Message.From.ID) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, ErrAccessDenied.Error())
		c.bot.Send(msg)

		return
	}

	pins, err := c.pinsParser.Parse()
	if err != nil {
		c.sendMessage(update.Message, fmt.Sprintf("parse error: %s", err))

		return
	}

	if err := c.repository.CreatePins(ctx, pins); err != nil {
		c.sendMessage(update.Message, fmt.Sprintf("save to db: %s", err))

		return
	}

	c.sendMessage(update.Message, fmt.Sprintf("successful. look it %s", CommandViewPins))
}
