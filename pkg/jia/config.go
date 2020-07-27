package jia

import (
	"os"
	"strconv"
)

type Config struct {
	BotToken          string
	ChannelID         string
	Port              int
	RedisURL          string
	VerificationToken string
}

func NewConfig() *Config {
	return &Config{
		BotToken:          getEnv("SLACK_CLIENT_BOT_TOKEN", ""),
		ChannelID:         getEnv("SLACK_CHANNEL_ID", ""),
		Port:              getEnvAsInt("PORT", 3000),
		RedisURL:          getEnv("REDIS_URL", "redis://localhost:6379/0"),
		VerificationToken: getEnv("SLACK_VERIFICATION_TOKEN", ""),
	}
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}

	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultValue
}
