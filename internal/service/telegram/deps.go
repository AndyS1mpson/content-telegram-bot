package telegram

import (
	"context"

	"content-telegram-bot/internal/models"
	"content-telegram-bot/internal/service/pin"
)

type pinsParser interface {
	Parse() ([]models.Pin, error)
}

type repository interface {
	GetPins(ctx context.Context, filter pin.PinFilter) ([]models.Pin, error)
	CreatePins(ctx context.Context, pins []models.Pin) error
}
