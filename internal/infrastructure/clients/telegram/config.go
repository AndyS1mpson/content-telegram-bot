package telegram

// Config конфигурация для клиента Telegram
type Config struct {
	Token string `env:"TELEGRAM_APITOKEN"`
}
