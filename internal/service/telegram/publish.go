package telegram

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"content-telegram-bot/internal/models"
	"content-telegram-bot/internal/utils/log"
)

const mediaGroupLimit = 10

// PublishHandler обрабатывает команду /publish <query> — публикует отобранные пины пачкой в канал.
func (c *TelegramClient) PublishHandler(ctx context.Context, update *tgbotapi.Update) {
	if !c.validateUser(update.Message.From.ID) {
		c.sendMessage(update.Message.Chat.ID, ErrAccessDenied.Error())
		return
	}

	query := strings.TrimSpace(update.Message.CommandArguments())
	if query == "" {
		c.sendMessage(update.Message.Chat.ID, ErrQueryRequired.Error())
		return
	}

	pins, err := c.pinService.GetSelected(ctx, c.defaultAccount.Channel, query)
	if err != nil {
		c.sendMessage(update.Message.Chat.ID, "Ошибка получения выбранных пинов: "+err.Error())
		return
	}

	if len(pins) == 0 {
		c.sendMessage(update.Message.Chat.ID, fmt.Sprintf("По теме \"%s\" нет отобранных пинов.", query))
		return
	}

	posted := make([]int64, 0, len(pins))
	for chunkStart := 0; chunkStart < len(pins); chunkStart += mediaGroupLimit {
		end := chunkStart + mediaGroupLimit
		if end > len(pins) {
			end = len(pins)
		}
		chunk := pins[chunkStart:end]

		media := buildMediaGroup(chunk)
		mg := tgbotapi.NewMediaGroup(c.defaultAccount.TelegramChatID, media)

		if _, err := c.bot.SendMediaGroup(mg); err != nil {
			log.Error(err, log.Data{"chunk_start": chunkStart})
			c.sendMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка публикации (опубликовано %d из %d): %s", len(posted), len(pins), err))
			if len(posted) > 0 {
				if markErr := c.pinService.MarkPosted(ctx, posted); markErr != nil {
					log.Error(markErr, log.Data{})
				}
			}
			return
		}

		for _, p := range chunk {
			posted = append(posted, p.ID)
		}
	}

	if err := c.pinService.MarkPosted(ctx, posted); err != nil {
		c.sendMessage(update.Message.Chat.ID, "Опубликовано, но не удалось обновить статусы: "+err.Error())
		return
	}

	c.sendMessage(update.Message.Chat.ID, fmt.Sprintf("Опубликовано %d пинов по теме \"%s\".", len(posted), query))
}

func buildMediaGroup(pins []models.Pin) []interface{} {
	media := make([]interface{}, 0, len(pins))
	for _, p := range pins {
		switch p.Type {
		case models.TypeVideo:
			media = append(media, tgbotapi.NewInputMediaVideo(tgbotapi.FileURL(p.URL)))
		default:
			media = append(media, tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(p.URL)))
		}
	}
	return media
}
