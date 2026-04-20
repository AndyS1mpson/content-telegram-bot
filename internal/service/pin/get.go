package pin

import (
	"context"

	"github.com/pkg/errors"

	"content-telegram-bot/internal/models"
)

// GetByID возвращает пин по его ID. Возвращает nil, nil если не найден.
func (s *Service) GetByID(ctx context.Context, id int64) (*models.Pin, error) {
	pins, err := s.repository.GetPins(ctx, Filter{IDs: []int64{id}})
	if err != nil {
		return nil, errors.Wrap(err, "get pin by id")
	}
	if len(pins) == 0 {
		return nil, nil
	}
	return &pins[0], nil
}
