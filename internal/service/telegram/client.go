package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type TelegramClient struct {
	pinsParser pinsParser
	repository      repository
	bot             *tgbotapi.BotAPI

	ownerID int64
}

func New(
	pinsParser pinsParser,
	repository repository,
	config Config,
) (*TelegramClient, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create telegram bot")
	}

	bot.Debug = false

	return &TelegramClient{
		pinsParser: pinsParser,
		repository:      repository,
		bot:             bot,
		ownerID:         config.BotOwnerID,
	}, nil
}

func (c *TelegramClient) validateUser(userID int64) bool {
	return userID == c.ownerID
}

// sendMessage отправляет сообщение в телеграм
func (c *TelegramClient) sendMessage(message *tgbotapi.Message, response string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ReplyToMessageID = message.MessageID
	c.bot.Send(msg)
}
