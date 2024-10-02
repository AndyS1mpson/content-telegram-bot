package pinterest

import (
	"github.com/pkg/errors"
	"github.com/playwright-community/playwright-go"

	"content-telegram-bot/internal/service/common"
)

const pinterestLoginURL = "https://ru.pinterest.com/login/"

type Parser struct {
	browser playwright.Browser

	tgChannel string
	login     string
	password  string
}

func New(browser playwright.Browser, config Config) *Parser {
	return &Parser{
		browser:   browser,
		tgChannel: config.TgChannel,
		login:     config.Login,
		password:  config.Password,
	}
}

// signIn попытка залогиниться на странице
func (p *Parser) signIn(page playwright.Page) error {
	if _, err := page.Goto(pinterestLoginURL); err != nil {
		return errors.Wrap(err, "go to login page")
	}
	if err := page.Locator(`input[name="id"]`).Fill(p.login); err != nil {
		return errors.Wrap(err, "fill login")
	}
	if err := page.Locator(`input[name="password"]`).Fill(p.password); err != nil {
		return errors.Wrap(err, "fill password")
	}
	if err := page.Locator(`button[type="submit"]`).Click(); err != nil {
		return errors.Wrap(err, "click button")
	}

	if err := page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:   playwright.LoadStateLoad,
		Timeout: playwright.Float(10000),
	}); err != nil {
		return errors.Wrap(err, "wait for login completed")
	}

	return nil
}

func (p *Parser) getNewPage() (playwright.Page, error) {
	return p.browser.NewPage(playwright.BrowserNewPageOptions{
		BaseURL:           playwright.String(common.BaseURL),
		JavaScriptEnabled: playwright.Bool(true),
		Viewport: &playwright.Size{
			Width:  1920,
			Height: 4000,
		},
		Screen: &playwright.Size{
			Width:  1920,
			Height: 4000,
		},
		UserAgent: playwright.String(common.GetUserAgent()),
		Locale:    playwright.String(common.Locale),
	})
}
