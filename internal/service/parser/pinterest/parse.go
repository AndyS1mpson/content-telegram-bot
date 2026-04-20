package pinterest

import (
	"regexp"
	"strconv"

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

// Parse получает изображения по заданной тематике (или с ленты рекомендаций, если query пустой).
func (p *Parser) Parse(account models.Account, query string) ([]models.Pin, error) {
	page, err := p.getNewPage()
	if err != nil {
		return nil, errors.Wrap(err, "create new page")
	}
	defer page.Close()

	if err := p.signIn(page, account); err != nil {
		return nil, errors.Wrap(err, "sign in")
	}

	if query != "" {
		if err := p.gotoSearch(page, query); err != nil {
			return nil, errors.Wrap(err, "go to search")
		}
	}

	pins, err := p.getPinsInfo(page, account, query)
	if err != nil {
		return nil, errors.Wrap(err, "parse images info")
	}

	return pins, nil
}

// getPinsInfo получение данных о пинах
func (p *Parser) getPinsInfo(page playwright.Page, account models.Account, query string) ([]models.Pin, error) {
	imagesLocator := page.Locator("img[srcset]")
	if err := imagesLocator.First().WaitFor(); err != nil {
		return nil, errors.Wrap(err, "wait for images locator")
	}

	rowPins, err := page.Evaluate(getPinInfoRowFunc)
	if err != nil {
		return nil, errors.Wrap(err, "get pins info")
	}

	rows, ok := rowPins.([]interface{})
	if !ok {
		return nil, errors.New("unexpected pins payload shape")
	}

	pins := make([]models.Pin, 0, len(rows))
	seen := make(map[int64]struct{}, len(rows))

	for _, raw := range rows {
		pinMap, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		idStr, _ := pinMap["id"].(string)
		rawURL, _ := pinMap["url"].(string)

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || rawURL == "" {
			continue
		}

		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}

		pins = append(pins, models.Pin{
			ID:      id,
			URL:     transformImageURL(rawURL),
			Type:    models.TypePin,
			Channel: account.Channel,
			Query:   query,
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
