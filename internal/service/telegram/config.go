package telegram

type Config struct {
	BotOwnerID int64 `yaml:"bot_owner_id"` // Идентификатор владельца бота
	Token string `yaml:"api_token"`
}
