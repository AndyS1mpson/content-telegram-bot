package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"content-telegram-bot/internal/utils/log"
)

// formatCallback формирует callback_data вида "<action>:<id>".
func formatCallback(action string, pinID int64) string {
	return fmt.Sprintf("%s:%d", action, pinID)
}

// parseCallback разбирает callback_data на action и pinID.
func parseCallback(data string) (action string, pinID int64, ok bool) {
	parts := strings.SplitN(data, ":", 2)
	if len(parts) != 2 {
		return "", 0, false
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", 0, false
	}
	return parts[0], id, true
}

// CallbackHandler обрабатывает нажатия на inline-кнопки под пинами.
func (c *TelegramClient) CallbackHandler(ctx context.Context, update *tgbotapi.Update) {
	cb := update.CallbackQuery

	if !c.validateUser(cb.From.ID) {
		c.answerCallback(cb.ID, ErrAccessDenied.Error())
		return
	}

	action, pinID, ok := parseCallback(cb.Data)
	if !ok {
		c.answerCallback(cb.ID, "Не удалось распознать действие")
		return
	}

	pin, err := c.pinService.GetByID(ctx, pinID)
	if err != nil || pin == nil {
		c.answerCallback(cb.ID, "Пин не найден")
		return
	}

	var markLabel string
	switch action {
	case callbackLike:
		if err := c.pinService.Select(ctx, pinID); err != nil {
			log.Error(err, log.Data{"pin_id": pinID})
			c.answerCallback(cb.ID, "Ошибка при сохранении выбора")
			return
		}
		markLabel = "✅ В публикации"
	case callbackDislike:
		markLabel = "❌ Отклонён"
	case callbackSkip:
		markLabel = "⏭ Пропущен"
	default:
		c.answerCallback(cb.ID, "Неизвестное действие")
		return
	}

	c.answerCallback(cb.ID, "")
	c.updateMessageAfterChoice(cb, pin.Query, markLabel)
	c.showNextPin(ctx, cb.Message.Chat.ID, pin.Channel, pin.Query)
}

func (c *TelegramClient) answerCallback(callbackID, text string) {
	cb := tgbotapi.NewCallback(callbackID, text)
	if _, err := c.bot.Request(cb); err != nil {
		log.Error(err, log.Data{"callback_id": callbackID})
	}
}

// updateMessageAfterChoice снимает клавиатуру и добавляет отметку выбора в caption.
func (c *TelegramClient) updateMessageAfterChoice(cb *tgbotapi.CallbackQuery, query, mark string) {
	caption := fmt.Sprintf("Тема: %s\n%s", query, mark)

	editCaption := tgbotapi.NewEditMessageCaption(cb.Message.Chat.ID, cb.Message.MessageID, caption)
	editCaption.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}}

	if _, err := c.bot.Request(editCaption); err != nil {
		log.Error(err, log.Data{"message_id": cb.Message.MessageID})
	}
}
