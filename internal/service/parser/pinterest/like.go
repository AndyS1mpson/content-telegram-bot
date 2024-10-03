package pinterest

import (
	"content-telegram-bot/internal/models"
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/playwright-community/playwright-go"
)

const pinURLTemplate = "https://ru.pinterest.com/pin/%v/"

// LikePins лайк пинов для обновления ленты рекомендаций
func (p *Parser) LikePins(ctx context.Context, pins []models.Pin, account models.Account) error {
	page, err := p.getNewPage()
	if err != nil {
		return errors.Wrap(err, "create new page")
	}
	defer page.Close()

	if err := p.signIn(page, account); err != nil {
		return errors.Wrap(err, "sign in")
	}

	for _, pin := range pins {
		if err := likePin(page, pin); err != nil {
			return errors.Wrap(err, "like pin")
		}
	}

	return nil
}

func likePin(page playwright.Page, pin models.Pin) error {
	_, err := page.Goto(fmt.Sprintf(pinURLTemplate, pin.ID))
	if err != nil {
		return errors.Wrap(err, "go to pin page")
	}

	if err := page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{State: playwright.LoadStateLoad}); err != nil {
		return errors.Wrap(err, "wait for pin page load state")
	}

	button := page.Locator("button", playwright.PageLocatorOptions{HasText: "Сохранить"})

	if err := button.WaitFor(); err != nil {
		return errors.Wrap(err, "wait for button load")
	}

	if err := button.Click(); err != nil {
		return errors.Wrap(err, "click save button")
	}

	return nil
}
