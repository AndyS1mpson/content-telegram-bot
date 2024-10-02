package models

// Pin информация о пине
type Pin struct {
	ID        int64    `json:"id"`
	ImageURL  string    `json:"image_url"`
	TgChannel string    `json:"tg_channel"`
	Status    PinStatus `json:"status"`
}

// PinStatus статус картинки
type PinStatus int64

var (
	PinStatusNew    PinStatus = 1 // новый спаршенный пин
	PinStatusViewed PinStatus = 2 // просмотренный пин
	PinStatusPosted PinStatus = 3 // запосченный пин
)
