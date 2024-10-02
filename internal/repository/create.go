package repository

import (
	"context"
	"errors"

	"content-telegram-bot/internal/models"

	sq "github.com/Masterminds/squirrel"
)

// CreatePins сохраняет информацию о пинах
func (r *Repository) CreatePins(ctx context.Context, pins []models.Pin) error {
	sql := sq.Insert(pinTableName).
		Columns("id", "image_url", "status", "tg_channel").
		PlaceholderFormat(sq.Dollar)

	for _, pin := range pins {
		sql = sql.Values(pin.ID, pin.ImageURL, pin.Status, pin.TgChannel)
	}

	raw, args, err := sql.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.ExecContext(ctx, raw, args...)
	if err != nil {
		return err
	}

	if add, _ := res.RowsAffected(); add != int64(len(pins)) {
		return errors.New("insert pinterest_pin rowsAffected error")
	}

	return nil
}
