package pin

import (
	"context"

	"content-telegram-bot/internal/models"
)

type parser interface {
	GetChannel() models.Channel
	Parse(account models.Account) ([]models.Pin, error)
}

type repository interface {
	GetPins(ctx context.Context, filter Filter) ([]models.Pin, error)
	CountPins(ctx context.Context, filter Filter) (int64, error)
	CreatePins(ctx context.Context, pins []models.Pin) error
	UpdateStatuses(ctx context.Context, ids []int64, newStatus models.PinStatus) error
}
