package pin

import (
	"context"

	"github.com/pkg/errors"

	"content-telegram-bot/internal/models"
)

// Select помечает пин как выбранный для публикации.
func (s *Service) Select(ctx context.Context, pinID int64) error {
	if err := s.repository.UpdateStatuses(ctx, []int64{pinID}, models.PinStatusSelected); err != nil {
		return errors.Wrap(err, "mark pin selected")
	}
	return nil
}
