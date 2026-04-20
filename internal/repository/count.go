package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"content-telegram-bot/internal/service/pin"
)

// CountPins посчитать количество пинов, удовлетворяющих фильтрам
func (r *Repository) CountPins(ctx context.Context, filter pin.Filter) (int64, error) {
	sql := sq.Select("COUNT(*) as total").From(pinTableName).PlaceholderFormat(sq.Dollar)

	sql = applyFilter(sql, filter)

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
