package storage

// Config настройки подключения к БД
type Config struct {
	User     string `yaml:"user"`
	DB       string `yaml:"db"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
}
