package pin

import (
	"context"

	"github.com/pkg/errors"

	"content-telegram-bot/internal/models"
)

// Parse парсинг пинов по теме и сохранение в БД
func (s *Service) Parse(ctx context.Context, account models.Account, query string) error {
	pins, err := s.parser.Parse(account, query)
	if err != nil {
		return errors.Wrap(err, "parse error")
	}

	if err := s.repository.CreatePins(ctx, pins); err != nil {
		return errors.Wrap(err, "save pins")
	}

	return nil
}
