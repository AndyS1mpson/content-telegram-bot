package models

// Pin информация о пине
type Pin struct {
	ID      int64     `json:"id"`      // идентификатор записи
	URL     string    `json:"url"`     // идентификатор контента
	Type    Type      `json:"type"`    // тип контента
	Status  PinStatus `json:"status"`  // статус
	Channel Channel   `json:"channel"` // канал для которого пин спаршен
}

// PinStatus статус картинки
type PinStatus int64

var (
	PinStatusNew    PinStatus = 1 // новый спаршенный объект
	PinStatusViewed PinStatus = 2 // просмотренный объект
	PinStatusPosted PinStatus = 3 // запосченный объект
)

type Channel string

var (
	ChannelWallpaper string = "Wall Paper"
)

// MediaType тип контента который будет поститься в тг канал
type Type string

var (
	TypePin   Type = "pin"
	TypeVideo Type = "video"
)
