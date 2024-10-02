package main

import (
	"context"
	"os"

	"github.com/pkg/errors"

	"content-telegram-bot/internal/config"
	"content-telegram-bot/internal/infrastructure/clients/browser"
	"content-telegram-bot/internal/service/pinterest"
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
	config, err := config.NewConfig()
	if err != nil {
		return failExitCode
	}

	br, close, err := browser.NewBrowser()
	defer close()
	if err != nil {
		log.Error(err, log.Data{})
		return failExitCode
	}

	parser := pinterest.New(br, pinterest.Config{
		Login:       config.PinterestParser[0].Login,
		Password:    config.PinterestParser[0].Password,
		ContentType: config.PinterestParser[0].ContentType,
	})

	_, err = parser.ParseFirstPage(context.Background())
	if err != nil {
		log.Error(errors.Wrap(err, "run parser"), log.Data{})
		return failExitCode
	}

	return successExitCode
}
