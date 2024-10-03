package pinterest

import (
	"regexp"

	"github.com/pkg/errors"
	"github.com/playwright-community/playwright-go"

	"content-telegram-bot/internal/models"
)

const getPinInfoRowFunc = `() => {
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

// Parse получает изображения с первой страницы ленты рекомендаций пинтереста
func (p *Parser) Parse(account models.Account) ([]models.Pin, error) {
	page, err := p.getNewPage()
	if err != nil {
		return nil, errors.Wrap(err, "create new page")
	}
	defer page.Close()

	if err := p.signIn(page, account); err != nil {
		return nil, errors.Wrap(err, "sign in")
	}

	pins, err := p.getPinsInfo(page, account)
	if err != nil {
		return nil, errors.Wrap(err, "parse images info")
	}

	return pins, nil
}

// getPinsInfo получение данных о пинах
func (p *Parser) getPinsInfo(page playwright.Page, account models.Account) ([]models.Pin, error) {
	imagesLocator := page.Locator("img[srcset]")
	if err := imagesLocator.First().WaitFor(); err != nil {
		return nil, errors.Wrap(err, "wait for images locator")
	}

	rowPins, err := page.Evaluate(getPinInfoRowFunc)
	if err != nil {
		return nil, errors.Wrap(err, "get pins info")
	}

	pins := make([]models.Pin, 0)

	for _, pin := range rowPins.([]interface{}) {
		pinMap := pin.(map[string]interface{})
		transformedURL := transformImageURL(pinMap["url"].(string))

		pins = append(pins, models.Pin{
			ID:      pinMap["id"].(int64),
			URL:     transformedURL,
			Channel: account.Channel,
			Status:  models.PinStatusNew,
		})
	}

	return pins, nil
}

func transformImageURL(url string) string {
	// Определяем регулярное выражение для поиска частей с размерами (например, '564x', '236x' и т.д.)
	pattern := regexp.MustCompile(`(https://i\.pinimg\.com/)\d+x/`)
	replacement := "${1}originals/"

	// Замена найденного паттерна на "originals"
	newURL := pattern.ReplaceAllString(url, replacement)
	return newURL
}
