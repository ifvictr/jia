package jia

import (
	"os"
)

type Config struct {
	BotToken  string
	ChannelId string
}

func NewConfig() *Config {
	return &Config{
		BotToken:  getEnv("SLACK_CLIENT_BOT_TOKEN", ""),
		ChannelId: getEnv("SLACK_CHANNEL_ID", ""),
	}
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}

	return defaultValue
}
