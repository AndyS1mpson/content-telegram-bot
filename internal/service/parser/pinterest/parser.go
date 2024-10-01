package pinterest

import (
	"context"
	"regexp"

	"github.com/pkg/errors"
	"github.com/playwright-community/playwright-go"

	"content-telegram-bot/internal/model"
	"content-telegram-bot/internal/service/parser/common"
)

const (
	pinterestLoginURL = "https://ru.pinterest.com/login/"

	getPinInfoRowFunc = `() => {
        const images = document.querySelectorAll('img[srcset]');
            const result = [];
        images.forEach(img => {
            const parent = img.closest('a[href*="/pin/"]');
            if (parent) {
                const pinId = parent.href.split("/pin/")[1].split("/")[0];
                result.push({id: pinId, url: img.src});
            }
        });
        return result.slice(0, 30);
        }`
)

type Parser struct {
	browser playwright.Browser

	contentType string
	login       string
	password    string
}

func New(browser playwright.Browser, config Config) *Parser {
	return &Parser{
		browser:     browser,
		contentType: config.ContentType,
		login:       config.Login,
		password:    config.Password,
	}
}

// ParseFirstPage получает изображения с первой страницы ленты рекомендаций пинтереста
func (p *Parser) ParseFirstPage(ctx context.Context) ([]model.PinterestPin, error) {
	page, err := p.getNewPage()
	if err != nil {
		return nil, errors.Wrap(err, "create new page")
	}
	defer page.Close()

	if err := p.signIn(page); err != nil {
		return nil, errors.Wrap(err, "sign in")
	}

	pins, err := p.getPinsInfo(page)
	if err != nil {
		return nil, errors.Wrap(err, "parse images info")
	}

	return pins, nil
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

// getPinsInfo получение данных о пинах
func (p *Parser) getPinsInfo(page playwright.Page) ([]model.PinterestPin, error) {
	imagesLocator := page.Locator("img[srcset]")
	if err := imagesLocator.First().WaitFor(); err != nil {
		return nil, errors.Wrap(err, "wait for images locator")
	}

	rowPins, err := page.Evaluate(getPinInfoRowFunc)
	if err != nil {
		return nil, errors.Wrap(err, "get pins info")
	}

	pins := make([]model.PinterestPin, 0)

	for _, pin := range rowPins.([]interface{}) {
		pinMap := pin.(map[string]interface{})
		transformedURL := transformImageURL(pinMap["url"].(string))

		pins = append(pins, model.PinterestPin{
			ID:       pinMap["id"].(string),
			ImageURL: transformedURL,
		})
	}

	return pins, nil
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

func transformImageURL(url string) string {
	// Определяем регулярное выражение для поиска частей с размерами (например, '564x', '236x' и т.д.)
	pattern := regexp.MustCompile(`(https://i\.pinimg\.com/)\d+x/`)
	replacement := "${1}originals/"

	// Замена найденного паттерна на "originals"
	newURL := pattern.ReplaceAllString(url, replacement)
	return newURL
}
