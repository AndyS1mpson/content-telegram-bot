package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"content-telegram-bot/internal/models"
	"content-telegram-bot/internal/utils/log"
)

type TelegramClient struct {
	bot        *tgbotapi.BotAPI
	pinService pinService

	accounts       map[models.Channel]models.Account
	defaultAccount models.Account

	ownerID int64
}

func New(
	bot *tgbotapi.BotAPI,
	pinService pinService,
	accounts map[models.Channel]models.Account,
	config Config,
) (*TelegramClient, error) {
	if len(accounts) == 0 {
		return nil, fmt.Errorf("at least one account must be configured")
	}

	var defaultAccount models.Account
	for _, acc := range accounts {
		defaultAccount = acc
		break
	}

	return &TelegramClient{
		bot:            bot,
		accounts:       accounts,
		defaultAccount: defaultAccount,
		pinService:     pinService,
		ownerID:        config.BotOwnerID,
	}, nil
}

func (c *TelegramClient) validateUser(userID int64) bool {
	return userID == c.ownerID
}

// sendMessage отправляет сообщение в телеграм
func (c *TelegramClient) sendMessage(chatID int64, response string) {
	msg := tgbotapi.NewMessage(chatID, response)
	if _, err := c.bot.Send(msg); err != nil {
		log.Error(err, log.Data{"chat_id": chatID, "text": response})
	}
}

// sendPinWithCheckboxes отправка пина с кнопками выбора
func (c *TelegramClient) sendPinWithCheckboxes(
	chatID int64,
	pin models.Pin,
	remainingCount int64,
) {
	likeCallback := formatCallback(callbackLike, pin.ID)
	dislikeCallback := formatCallback(callbackDislike, pin.ID)
	skipCallback := formatCallback(callbackSkip, pin.ID)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❤️ В публикацию", likeCallback),
			tgbotapi.NewInlineKeyboardButtonData("👎 Отклонить", dislikeCallback),
			tgbotapi.NewInlineKeyboardButtonData("⏭ Пропустить", skipCallback),
		),
	)

	caption := fmt.Sprintf("Тема: %s\nОсталось неотсмотренных: %d", pin.Query, remainingCount)

	if err := c.sendContent(chatID, pin, caption, &keyboard); err != nil {
		c.sendMessage(chatID, err.Error())
	}
}

// sendContent универсальная функция для отправки картинки или видео
func (c *TelegramClient) sendContent(
	chatID int64,
	pin models.Pin,
	caption string,
	replyMarkup *tgbotapi.InlineKeyboardMarkup,
) error {
	var msg tgbotapi.Chattable

	switch pin.Type {
	case models.TypeVideo:
		video := tgbotapi.NewVideo(chatID, tgbotapi.FileURL(pin.URL))
		video.Caption = caption
		if replyMarkup != nil {
			video.ReplyMarkup = replyMarkup
		}
		msg = video
	default:
		photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(pin.URL))
		photo.Caption = caption
		if replyMarkup != nil {
			photo.ReplyMarkup = replyMarkup
		}
		msg = photo
	}

	if _, err := c.bot.Send(msg); err != nil {
		return err
	}

	return nil
}

// RegisterHandlers регистрация обработчиков сообщений и callback-запросов
func (c *TelegramClient) RegisterHandlers(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.bot.GetUpdatesChan(u)

	for update := range updates {
		switch {
		case update.CallbackQuery != nil:
			c.CallbackHandler(ctx, &update)
		case update.Message != nil:
			c.handleMessage(ctx, &update)
		}
	}
}

func (c *TelegramClient) handleMessage(ctx context.Context, update *tgbotapi.Update) {
	switch Command("/" + update.Message.Command()) {
	case CommandStart:
		c.StartHandler(ctx, update)
	case CommandCollect:
		c.CollectHandler(ctx, update)
	case CommandView:
		c.ViewHandler(ctx, update)
	case CommandPublish:
		c.PublishHandler(ctx, update)
	default:
		c.sendMessage(update.Message.Chat.ID, "Unknown command. Try /start")
	}
}
