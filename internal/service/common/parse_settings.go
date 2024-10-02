package common

import "math/rand"

const (
	Locale = "ru-RU"
	BaseURL = "https://www.google.com/"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
}

// GetUserAgent получение случайного User-Agent
func GetUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}
