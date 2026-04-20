package pin

import (
	"context"

	"github.com/pkg/errors"

	"content-telegram-bot/internal/models"
)

// GetSelected возвращает все пины, отмеченные как Selected, для указанного канала и темы.
func (s *Service) GetSelected(ctx context.Context, channel models.Channel, query string) ([]models.Pin, error) {
	pins, err := s.repository.GetPins(ctx, Filter{
		Statuses: []models.PinStatus{models.PinStatusSelected},
		Channels: []models.Channel{channel},
		Query:    &query,
	})
	if err != nil {
		return nil, errors.Wrap(err, "get selected pins")
	}
	return pins, nil
}

// MarkPosted переводит указанные пины в статус Posted.
func (s *Service) MarkPosted(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.repository.UpdateStatuses(ctx, ids, models.PinStatusPosted); err != nil {
		return errors.Wrap(err, "mark pins posted")
	}
	return nil
}
