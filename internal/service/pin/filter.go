package pin

import (
	"content-telegram-bot/internal/models"
	"content-telegram-bot/internal/utils/sql"
)

// PinFilter фильтры для поиска пинов
type PinFilter struct {
	IDs        []int64
	Statuses   []models.PinStatus
	TgChannels []string
	PageQuery  *sql.PageQuery
}
