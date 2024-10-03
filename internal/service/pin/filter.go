package pin

import (
	"content-telegram-bot/internal/models"
)

// Filter фильтры для поиска контента
type Filter struct {
	IDs      []int64
	Statuses []models.PinStatus
	Types    []models.Type
	Channels []models.Channel
	Limit    *int64
}
