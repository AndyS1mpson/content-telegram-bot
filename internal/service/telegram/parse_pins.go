package telegram

import (
	"content-telegram-bot/internal/models"
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

	account, ok := c.accounts[models.Channel(update.Message.Text)]
	if !ok {
		c.sendMessage(update.Message.Chat.ID, ErrIncorrectAction.Error())
		return
	}

	if err := c.pinService.Parse(ctx, account); err != nil {
		c.sendMessage(update.Message.Chat.ID, fmt.Sprintf("parsing error: %s", err))
		return
	}

	c.sendMessage(update.Message.Chat.ID, fmt.Sprintf("successful. look it %s", CommandViewPins))
}
