package telegram

import "errors"

var (
	CommandParsePinterest Command = "/parse_pinterest" // Запуск парсинга пинов из pinterest'а
	CommandViewPins       Command = "/view_pins"       // Просмотр пинов
)

// Команда для телеграм бота
type Command string

var (
	ErrAccessDenied = errors.New("you do not have access")
	ErrIncorrectAction = errors.New("incorrect action")
)
