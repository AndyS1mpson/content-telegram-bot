package main

import (
	"os"

	"github.com/pkg/errors"

	"content-telegram-bot/internal/config"
	"content-telegram-bot/internal/infrastructure/clients/browser"
	"content-telegram-bot/internal/models"
	"content-telegram-bot/internal/service/parser/pinterest"
	"content-telegram-bot/internal/utils/log"
	"content-telegram-bot/internal/utils/slices"
)

const (
	successExitCode = 0
	failExitCode    = 1
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	cfg, err := config.NewConfig()
	if err != nil {
		return failExitCode
	}

	br, close, err := browser.NewBrowser()
	defer close()
	if err != nil {
		log.Error(err, log.Data{})
		return failExitCode
	}

	parser := pinterest.New(br)

	accounts := slices.Map(cfg.Accounts, func(config config.AccountConfig) models.Account {
		return models.Account{
			Channel:  models.Channel(config.Channel),
			Login:    config.Login,
			Password: config.Password,
		}
	})

	_, err = parser.Parse(accounts[0])
	if err != nil {
		log.Error(errors.Wrap(err, "run parser"), log.Data{})
		return failExitCode
	}

	return successExitCode
}
