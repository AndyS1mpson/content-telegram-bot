package repository

import (
	"context"

	"content-telegram-bot/internal/models"
	"content-telegram-bot/internal/service/pin"
	"content-telegram-bot/internal/utils/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

// GetPins поиск пинов по фильтрам
func (r *Repository) GetPins(ctx context.Context, filter pin.PinFilter) ([]models.Pin, error) {
	query := sq.Select("*").From(pinTableName).PlaceholderFormat(sq.Dollar)

	if filter.IDs != nil {
		query = query.Where(sq.Eq{"id": filter.IDs})
	}

	if filter.Statuses != nil {
		query = query.Where(sq.Eq{"status": filter.Statuses})
	}

	if filter.TgChannels != nil {
		query = query.Where(sq.Eq{"tg_channel": filter.TgChannels})
	}

	if filter.PageQuery != nil {
		query = sql.AddPaging(query, *filter.PageQuery, "id")
	}

	raw, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var entities []PinterestPin

	if err = r.db.SelectContext(ctx, &entities, raw, args...); err != nil {
		return nil, errors.Wrap(err, "find pins")
	}

	pins := make([]models.Pin, 0, len(entities))
	for _, entity := range entities {
		pins = append(pins, models.Pin{
			ID:        entity.ID,
			ImageURL:  entity.ImageURL,
			TgChannel: entity.TgChannel,
			Status:    models.PinStatus(entity.Status),
		})
	}

	return pins, nil
}
