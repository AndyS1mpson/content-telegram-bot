package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"content-telegram-bot/internal/models"
)

type TelegramClient struct {
	bot        *tgbotapi.BotAPI
	pinService pinService

	accounts map[models.Channel]models.Account

	ownerID int64
}

func New(
	bot *tgbotapi.BotAPI,
	pinService pinService,
	accounts map[models.Channel]models.Account,
	config Config,
) (*TelegramClient, error) {
	return &TelegramClient{
		bot:        bot,
		accounts:   accounts,
		pinService: pinService,
		ownerID:    config.BotOwnerID,
	}, nil
}

func (c *TelegramClient) validateUser(userID int64) bool {
	return userID == c.ownerID
}

// sendMessage отправляет сообщение в телеграм
func (c *TelegramClient) sendMessage(chatID int64, response string) {
	msg := tgbotapi.NewMessage(chatID, response)
	c.bot.Send(msg)
}

// sendPinWithCheckboxes отправка пина с командами
func (c *TelegramClient) sendPinWithCheckboxes(
	chatID int64,
	pin models.Pin,
	unwatchedPinsCount int64,
) {
	likeCallback := fmt.Sprintf("like_%d_%s_%s", pin.ID, pin.Type, pin.Channel)
	dislikeCallback := fmt.Sprintf("dislike_%d_%s_%s", pin.ID, pin.Type, pin.Channel)
	skipCallback := fmt.Sprintf("skip_%d_%s_%s", pin.ID, pin.Type, pin.Channel)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❤️", likeCallback),
			tgbotapi.NewInlineKeyboardButtonData("👎", dislikeCallback),
			tgbotapi.NewInlineKeyboardButtonData("Не хочу больше смотреть", skipCallback),
		),
	)

	caption := fmt.Sprintf("Осталось еще %d не просмотренных %s", unwatchedPinsCount, pin.Type)

	if err := c.sendContent(chatID, pin, caption, &keyboard); err != nil {
		c.sendMessage(chatID, err.Error())
	}

}

// sendImage универсальная функция для отправки контента
func (c *TelegramClient) sendContent(
	chatID int64,
	pin models.Pin,
	caption string,
	replyMarkup *tgbotapi.InlineKeyboardMarkup,
) error {
	var msg tgbotapi.Chattable

	switch pin.Type {
	case models.TypePin:
		msg = tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(pin.URL))
		msg.(*tgbotapi.PhotoConfig).Caption = caption
	case models.TypeVideo:
		msg = tgbotapi.NewVideo(chatID, tgbotapi.FileURL(pin.URL))
		msg.(*tgbotapi.VideoConfig).Caption = caption
	}

	// Если передана клавиатура, добавляем её к сообщению
	if replyMarkup != nil {
		switch m := msg.(type) {
		case *tgbotapi.PhotoConfig:
			m.ReplyMarkup = replyMarkup
		case *tgbotapi.VideoConfig:
			m.ReplyMarkup = replyMarkup
		}
	}

	_, err := c.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// RegisterHandlers регистрация обработчиков
func (c *TelegramClient) RegisterHandlers(ctx context.Context) {
	// Запуск polling для обработки обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch Command(update.Message.Command()) {
		case CommandStart:
			c.StartHandler(ctx, &update)
		case CommandViewPins:
			c.ViewNewPinHandler(ctx, &update)
		case CommandParsePinterest:
			c.ParsePinsHandler(ctx, &update)
		default:
			// Неизвестная команда
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
			c.bot.Send(msg)
		}
	}
}
