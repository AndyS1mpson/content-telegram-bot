package telegram

import (
	"context"
	"strings"

	"github.com/AlekSi/pointer"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"content-telegram-bot/internal/models"
	"content-telegram-bot/internal/service/pin"
)

// ViewHandler обрабатывает команду /view <query> — показывает следующий неотсмотренный пин.
func (c *TelegramClient) ViewHandler(ctx context.Context, update *tgbotapi.Update) {
	if !c.validateUser(update.Message.From.ID) {
		c.sendMessage(update.Message.Chat.ID, ErrAccessDenied.Error())
		return
	}

	query := strings.TrimSpace(update.Message.CommandArguments())
	if query == "" {
		c.sendMessage(update.Message.Chat.ID, ErrQueryRequired.Error())
		return
	}

	c.showNextPin(ctx, update.Message.Chat.ID, c.defaultAccount.Channel, query)
}

// showNextPin выдаёт следующий пин со статусом New. Если таких нет —
// автоматически дозапускает парсинг по указанной теме.
func (c *TelegramClient) showNextPin(ctx context.Context, chatID int64, channel models.Channel, query string) {
	pins, count, err := c.fetchNextPin(ctx, channel, query)
	if err != nil {
		c.sendMessage(chatID, err.Error())
		return
	}

	if len(pins) > 0 {
		c.sendPinWithCheckboxes(chatID, pins[0], count-1)
		return
	}

	c.sendMessage(chatID, "Новых пинов нет, собираю ещё…")

	if err := c.pinService.Parse(ctx, c.defaultAccount, query); err != nil {
		c.sendMessage(chatID, "Ошибка парсинга: "+err.Error())
		return
	}

	pins, count, err = c.fetchNextPin(ctx, channel, query)
	if err != nil {
		c.sendMessage(chatID, err.Error())
		return
	}

	if len(pins) == 0 {
		c.sendMessage(chatID, "По теме \""+query+"\" больше нечего показать.")
		return
	}

	c.sendPinWithCheckboxes(chatID, pins[0], count-1)
}

func (c *TelegramClient) fetchNextPin(ctx context.Context, channel models.Channel, query string) ([]models.Pin, int64, error) {
	return c.pinService.GetPinsForView(ctx, pin.Filter{
		Statuses: []models.PinStatus{models.PinStatusNew},
		Channels: []models.Channel{channel},
		Query:    &query,
		Limit:    pointer.ToInt64(1),
	})
}
