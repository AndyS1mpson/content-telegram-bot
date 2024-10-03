package pinterest

import (
	"github.com/pkg/errors"
	"github.com/playwright-community/playwright-go"

	"content-telegram-bot/internal/models"
	"content-telegram-bot/internal/service/common"
)

const pinterestLoginURL = "https://ru.pinterest.com/login/"

type Parser struct {
	browser playwright.Browser
}

func New(browser playwright.Browser) *Parser {
	return &Parser{
		browser: browser,
	}
}

// signIn попытка залогиниться на странице
func (p *Parser) signIn(page playwright.Page, account models.Account) error {
	if _, err := page.Goto(pinterestLoginURL); err != nil {
		return errors.Wrap(err, "go to login page")
	}
	if err := page.Locator(`input[name="id"]`).Fill(account.Login); err != nil {
		return errors.Wrap(err, "fill login")
	}
	if err := page.Locator(`input[name="password"]`).Fill(account.Password); err != nil {
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
