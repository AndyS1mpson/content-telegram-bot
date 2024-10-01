package config

import (
	"os"
	"path/filepath"
	"content-telegram-bot/internal/service/parser/pinterest"

	"gopkg.in/yaml.v2"
)

var configFileName = "config.yaml"

// AppConfig конфигурация приложения
type AppConfig struct {
	PinterestParser []pinterest.Config `yaml:"pinterest"` // конфигурации аккаунтов, разделенных по типу контента в рекомендациях
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
