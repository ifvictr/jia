package main

import (
	"fmt"
	"log"

	"github.com/ifvictr/fig/pkg/fig"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	config := fig.NewConfig()
	api := slack.New(config.BotToken)

	api.PostMessage(config.ChannelId, slack.MsgOptionText("Hello, world!", false))
	fmt.Println("Hello, world!")
}
