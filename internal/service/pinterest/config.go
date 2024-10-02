package pinterest

// Config конфигурация парсера
type Config struct {
	ContentType string `yaml:"content_type"` // Под разные тематики разные аккаунты чтобы настраивать ленту рекомендаций
	Login       string `yaml:"login"`
	Password    string `yaml:"password"`
}
