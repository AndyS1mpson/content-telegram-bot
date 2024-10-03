package telegram

import (
	"context"

	"github.com/AlekSi/pointer"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"content-telegram-bot/internal/models"
	"content-telegram-bot/internal/service/pin"
)

const noNewPins = "Новых пинов нет"

// ViewNewPinHandler обработчик вызова команды просмотра нового пина
func (c *TelegramClient) ViewNewPinHandler(ctx context.Context, update *tgbotapi.Update) {
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

	pins, count, err := c.pinService.GetPinsForView(ctx, pin.Filter{
		Statuses: []models.PinStatus{models.PinStatusNew},
		Channels: []models.Channel{account.Channel},
		Limit:    pointer.ToInt64(1),
	})
	if err != nil {
		c.sendMessage(update.Message.Chat.ID, err.Error())
	}

	if len(pins) == 0 {
		c.sendMessage(update.Message.Chat.ID, noNewPins)
		return
	}

	c.sendPinWithCheckboxes(update.Message.Chat.ID, pins[0], count - 1)
}
