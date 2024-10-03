package telegram

import (
	"content-telegram-bot/internal/models"
	"content-telegram-bot/internal/service/pin"
	"context"
)

type pinService interface {
	Parse(ctx context.Context, account models.Account) error
	GetPinsForView(ctx context.Context, filter pin.Filter) ([]models.Pin, int64, error)
}
