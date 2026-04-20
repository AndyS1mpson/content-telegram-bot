package telegram

import (
	"context"

	"content-telegram-bot/internal/models"
	"content-telegram-bot/internal/service/pin"
)

type pinService interface {
	Parse(ctx context.Context, account models.Account, query string) error
	GetPinsForView(ctx context.Context, filter pin.Filter) ([]models.Pin, int64, error)
	GetByID(ctx context.Context, id int64) (*models.Pin, error)
	Select(ctx context.Context, pinID int64) error
	GetSelected(ctx context.Context, channel models.Channel, query string) ([]models.Pin, error)
	MarkPosted(ctx context.Context, ids []int64) error
}
