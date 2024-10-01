package browser

import (
	"github.com/pkg/errors"
	"github.com/playwright-community/playwright-go"
)

// NewBrowser запуск нового браузера
func NewBrowser() (playwright.Browser, func() error, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, nil, errors.Wrap(err, "run playwright")
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "run browser")
	}

	return browser, func() error {
		if err := browser.Close(); err != nil {
			return err
		}
		return pw.Stop()
	}, nil
}
