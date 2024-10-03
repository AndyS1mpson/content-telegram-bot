package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"content-telegram-bot/internal/infrastructure/storage"
	"content-telegram-bot/internal/service/telegram"
)

var configFileName = "config.yaml"

// AccountConfig связь аккаунта с данными для парсинга и канала в который будет поститься контент
type AccountConfig struct {
	Channel  string `yaml:"channel"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
}

// AppConfig конфигурация приложения
type AppConfig struct {
	Telegram telegram.Config `yaml:"telegram"`
	Accounts []AccountConfig `yaml:"accounts"`
	Database storage.Config  `yaml:"database"`
}

// NewConfig возвращает декодированную конфигурацию приложения
func NewConfig() (*AppConfig, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(rootDir, configFileName)

	config := &AppConfig{}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
