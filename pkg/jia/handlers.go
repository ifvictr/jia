package jia

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func HandleInnerEvent(slackClient *slack.Client, innerEvent *slackevents.EventsAPIInnerEvent) {
	switch e := innerEvent.Data.(type) {
	case *slackevents.MessageEvent:
		onMessage(slackClient, e)
		break
	}
}

func onMessage(slackClient *slack.Client, event *slackevents.MessageEvent) {
	// Ignore messages that aren't in the target channel, or are non-user messages
	if event.Channel != jiaConfig.ChannelID || event.User == "USLACKBOT" || event.User == "" {
		return
	}

	// Attempt to extract a positive number from the beginning of a string
	countPattern := regexp.MustCompile(`^\d+`)
	matchedNumber, err := strconv.Atoi(countPattern.FindString(event.Text))
	log.Println(matchedNumber)

	// Ignore messages that don't have numbers.
	if err != nil {
		log.Println("Failed to retrieve number, skipping…")
		return
	}

	// Reject if sender also sent the previous number.
	lastSenderID, err := redisClient.Get("last_sender_id").Result()
	if err != nil {
		log.Println("Failed to retrieve the last sender")
		return
	}
	if event.User == lastSenderID {
		slackClient.AddReaction("bangbang", slack.ItemRef{
			Channel:   event.Channel,
			Timestamp: event.TimeStamp,
		})
		slackClient.PostEphemeral(event.Channel, event.User, slack.MsgOptionText(
			"You counted consecutively! That’s not allowed.", false))
		return
	}

	// Ignore numbers that aren't in order.
	lastValidNumberStr, err := redisClient.Get("last_valid_number").Result()
	if err != nil {
		log.Println("Failed to retrieve the last valid number")
		return
	}
	lastValidNumber, err := strconv.Atoi(lastValidNumberStr)
	if err != nil {
		log.Println("Failed to convert the last valid number to an integer")
		return
	}
	if matchedNumber != lastValidNumber+1 {
		slackClient.AddReaction("bangbang", slack.ItemRef{
			Channel:   event.Channel,
			Timestamp: event.TimeStamp,
		})
		slackClient.PostEphemeral(event.Channel, event.User, slack.MsgOptionText(
			fmt.Sprintf("You counted incorrectly! The next valid number is supposed to be *%d*.", lastValidNumber+1), false))
		return
	}

	// Finally!
	redisClient.Set("last_valid_number", matchedNumber, 0)
	redisClient.Set("last_sender_id", event.User, 0)
}
