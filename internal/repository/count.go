package repository

import (
	"context"

	"content-telegram-bot/internal/service/pin"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

// CountPins посчитать количество пинов, удовлетсоряющих фильтрам
func (r *Repository) CountPins(ctx context.Context, filter pin.Filter) (int64, error) {
	sql := sq.Select("COUNT(*) as total").From(pinTableName).PlaceholderFormat(sq.Dollar)

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
		return 0, err
	}

	var total int64

	if err = r.db.GetContext(ctx, &total, raw, args...); err != nil {
		return 0, errors.Wrap(err, "count pins")
	}

	return total, nil
}
