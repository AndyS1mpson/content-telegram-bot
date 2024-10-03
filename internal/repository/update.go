package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"content-telegram-bot/internal/models"
)

// UpdateStatuses обновление статусов пачки пинов
func (r *Repository) UpdateStatuses(ctx context.Context, pinIDs []int64, newStatus models.PinStatus) error {
	sql := sq.Update(pinTableName).SetMap(map[string]interface{}{
		"status": newStatus,
	}).Where(sq.Eq{"id": pinIDs}).PlaceholderFormat(sq.Dollar)

	raw, args, err := sql.ToSql()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, raw, args...); err != nil {
		return errors.Wrap(err, "update statuses")
	}

	return nil
}
