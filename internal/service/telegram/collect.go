package telegram

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CollectHandler обрабатывает команду /collect <query> — запускает парсинг Pinterest по теме.
func (c *TelegramClient) CollectHandler(ctx context.Context, update *tgbotapi.Update) {
	if !c.validateUser(update.Message.From.ID) {
		c.sendMessage(update.Message.Chat.ID, ErrAccessDenied.Error())
		return
	}

	query := strings.TrimSpace(update.Message.CommandArguments())
	if query == "" {
		c.sendMessage(update.Message.Chat.ID, ErrQueryRequired.Error())
		return
	}

	c.sendMessage(update.Message.Chat.ID, fmt.Sprintf("Собираю пины по теме: %s…", query))

	if err := c.pinService.Parse(ctx, c.defaultAccount, query); err != nil {
		c.sendMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка парсинга: %s", err))
		return
	}

	c.sendMessage(update.Message.Chat.ID, fmt.Sprintf("Готово. Смотри: /view %s", query))
}
