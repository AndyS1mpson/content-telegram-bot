package pin

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"content-telegram-bot/internal/models"
	"content-telegram-bot/internal/utils/slices"
)

// GetPinsForView получение пинов из хранилища и обновление статуса на "просмотрен"
func (s *Service) GetPinsForView(ctx context.Context, filter Filter) ([]models.Pin, int64, error) {
	var (
		count int64
		pins  []models.Pin
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		count, err = s.repository.CountPins(ctx, filter)
		if err != nil {
			return errors.Wrap(err, "get pins count")
		}

		return nil
	})

	g.Go(func() error {
		var err error
		pins, err = s.repository.GetPins(ctx, filter)
		if err != nil {
			return errors.Wrap(err, "get new content")
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	if err := s.repository.UpdateStatuses(
		ctx,
		slices.Map(pins, func(pin models.Pin) int64 { return pin.ID }),
		models.PinStatusViewed,
	); err != nil {
		return nil, 0, errors.Wrap(err, "update statuses for viewed content")
	}

	return pins, count, nil
}
