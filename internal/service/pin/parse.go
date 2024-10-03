package pin

import (
	"content-telegram-bot/internal/models"
	"context"

	"github.com/pkg/errors"
)

// Parse парсинг пинов и сохранение в БД
func (s *Service) Parse(ctx context.Context, account models.Account) error {
	pins, err := s.parser.Parse(account)
	if err != nil {
		return errors.Wrap(err, "parse error")
	}

	if err := s.repository.CreatePins(ctx, pins); err != nil {
		return errors.Wrap(err, "save pins")
	}

	return nil
}
