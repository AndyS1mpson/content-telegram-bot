package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"content-telegram-bot/internal/models"
	"content-telegram-bot/internal/service/pin"
)

// GetPins поиск пинов по фильтрам
func (r *Repository) GetPins(ctx context.Context, filter pin.Filter) ([]models.Pin, error) {
	sql := sq.Select("*").From(pinTableName).PlaceholderFormat(sq.Dollar)

	if filter.IDs != nil {
		sql = sql.Where(sq.Eq{"id": filter.IDs})
	}

	if filter.Statuses != nil {
		sql = sql.Where(sq.Eq{"status": filter.Statuses})
	}

	if filter.Channels != nil {
		sql = sql.Where(sq.Eq{"channel": filter.Channels})
	}

	if filter.Limit != nil {
		sql = sql.Limit(uint64(*filter.Limit))
	}

	if filter.Types != nil {
		sql = sql.Where(sq.Eq{"type": filter.Types})
	}

	raw, args, err := sql.ToSql()
	if err != nil {
		return nil, err
	}

	var entities []Pin

	if err = r.db.SelectContext(ctx, &entities, raw, args...); err != nil {
		return nil, errors.Wrap(err, "find pins")
	}

	pins := make([]models.Pin, 0, len(entities))
	for _, entity := range entities {
		pins = append(pins, models.Pin{
			ID:      entity.ID,
			URL:     entity.URL,
			Channel: models.Channel(entity.Channel),
			Status:  models.PinStatus(entity.Status),
		})
	}

	return pins, nil
}
