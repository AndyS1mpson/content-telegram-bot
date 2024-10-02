package pinterest

// Config конфигурация парсера
type Config struct {
	TgChannel string `yaml:"tg_channel"` // Под разные каналы (1 канал == 1 тематика) разные аккаунты чтобы настраивать ленту рекомендаций
	Login     string `yaml:"login"`
	Password  string `yaml:"password"`
}
