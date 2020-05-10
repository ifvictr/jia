package fig

import (
	"os"
)

type Config struct {
	BotToken            string
	ClientId            string
	ClientSecret        string
	ClientSigningSecret string
	ChannelId           string
}

func NewConfig() *Config {
	return &Config{
		BotToken:            getEnv("BOT_TOKEN", ""),
		ClientId:            getEnv("CLIENT_ID", ""),
		ClientSecret:        getEnv("CLIENT_SECRET", ""),
		ClientSigningSecret: getEnv("CLIENT_SIGNING_SECRET", ""),
		ChannelId:           getEnv("CHANNEL_ID", ""),
	}
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}

	return defaultValue
}
