package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"content-telegram-bot/internal/models"
)

// CreatePins сохраняет информацию о пинах. Дубликаты по (id, channel) пропускаются.
func (r *Repository) CreatePins(ctx context.Context, pins []models.Pin) error {
	if len(pins) == 0 {
		return nil
	}

	sql := sq.Insert(pinTableName).
		Columns("id", "url", "status", "type", "channel", "query").
		PlaceholderFormat(sq.Dollar).
		Suffix("ON CONFLICT (id, channel) DO NOTHING")

	for _, pin := range pins {
		sql = sql.Values(pin.ID, pin.URL, pin.Status, pin.Type, pin.Channel, pin.Query)
	}

	raw, args, err := sql.ToSql()
	if err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx, raw, args...); err != nil {
		return errors.Wrap(err, "insert pins")
	}

	return nil
}
