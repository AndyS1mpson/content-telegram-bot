package main

import (
	"context"
	"os"

	// "github.com/go-co-op/gocron"

	"content-telegram-bot/internal/config"
	"content-telegram-bot/internal/infrastructure/clients/browser"
	tgClient "content-telegram-bot/internal/infrastructure/clients/telegram"
	"content-telegram-bot/internal/infrastructure/storage"
	"content-telegram-bot/internal/models"
	"content-telegram-bot/internal/repository"
	"content-telegram-bot/internal/service/parser/pinterest"
	"content-telegram-bot/internal/service/pin"
	"content-telegram-bot/internal/service/telegram"
	"content-telegram-bot/internal/utils/log"
)

const (
	successExitCode = 0
	failExitCode    = 1
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	ctx := context.Background()

	// Config
	cfg, err := config.NewConfig()
	if err != nil {
		return failExitCode
	}

	// DB
	db := storage.New(cfg.Database)
	defer db.Close()

	// Repository
	pinRepository := repository.New(db)

	// Browser
	br, close, err := browser.NewBrowser()
	defer close()
	if err != nil {
		log.Error(err, log.Data{})
		return failExitCode
	}

	// Parser
	parser := pinterest.New(br)

	// Service
	pinService := pin.NewService(parser, pinRepository)

	// Client
	bot, err := tgClient.New(cfg.Telegram.Token)
	if err != nil {
		log.Error(err, log.Data{})
		return failExitCode
	}

	tgWrapper, err := telegram.New(bot, pinService, getAccounts(cfg.Accounts), cfg.Telegram)
	if err != nil {
		log.Error(err, log.Data{})
		return failExitCode
	}

	tgWrapper.RegisterHandlers(ctx)

	return successExitCode
}

func getAccounts(cfg []config.AccountConfig) map[models.Channel]models.Account {
	result := make(map[models.Channel]models.Account, len(cfg))

	for _, account := range cfg {
		result[models.Channel(account.Channel)] = models.Account{
			Channel:  models.Channel(account.Channel),
			Login:    account.Login,
			Password: account.Password,
		}
	}

	return result
}
