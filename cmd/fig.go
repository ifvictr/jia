package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/ifvictr/fig/pkg/fig"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func init() {
	godotenv.Load()
}

func main() {
	fmt.Println("Starting Fig…")
	config := fig.NewConfig()

	api := slack.New(config.BotToken)
	rtm := api.NewRTM()

	lastValidNumber := 0
	var lastSender string

	go rtm.ManageConnection()

	for event := range rtm.IncomingEvents {
		switch ev := event.Data.(type) {
		case *slack.MessageEvent:
			// Ignore messages that aren't in the target channel, or are non-user messages
			if ev.Channel != config.ChannelId || ev.User == "USLACKBOT" || ev.User == "" {
				continue
			}

			// Attempt to extract a positive number from the beginning of a string
			countPattern := regexp.MustCompile(`^\d+`)
			matchedNumber, err := strconv.Atoi(countPattern.FindString(ev.Text))
			log.Println(matchedNumber)

			// Ignore messages that don't have numbers
			if err != nil {
				log.Println("Failed to retrieve number, skipping…")
				continue
			}

			// Reject if sender also sent the previous message
			if ev.User == lastSender {
				api.AddReaction("bangbang", slack.ItemRef{
					Channel:   ev.Channel,
					Timestamp: ev.Timestamp,
				})
				api.PostEphemeral(ev.Channel, ev.User, slack.MsgOptionText("You counted consecutively! That’s not allowed.", false))
				continue
			}

			// Ignore numbers that aren't in order
			if matchedNumber != lastValidNumber+1 {
				api.AddReaction("bangbang", slack.ItemRef{
					Channel:   ev.Channel,
					Timestamp: ev.Timestamp,
				})
				api.PostEphemeral(ev.Channel, ev.User, slack.MsgOptionText("You counted incorrectly! The next valid number is "+strconv.Itoa(lastValidNumber+1), false))
				continue
			}

			// Finally!
			lastValidNumber = matchedNumber
			lastSender = ev.User
			api.AddReaction("+1", slack.ItemRef{
				Channel:   ev.Channel,
				Timestamp: ev.Timestamp,
			})
		}
	}
}
