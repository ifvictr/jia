package jia

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func HandleInnerEvent(slackClient *slack.Client, innerEvent *slackevents.EventsAPIInnerEvent) {
	switch e := innerEvent.Data.(type) {
	case *slackevents.MessageEvent:
		onMessage(slackClient, e)
	}
}

func onMessage(slackClient *slack.Client, event *slackevents.MessageEvent) {
	// Ignore messages that aren't in the target channel, or are non-user messages.
	if event.Channel != jiaConfig.ChannelID || event.User == "USLACKBOT" || event.User == "" {
		return
	}

	// Ignore threaded messages.
	if event.ThreadTimeStamp != "" {
		return
	}

	// Attempt to extract a positive number at the start of a string.
	countPattern := regexp.MustCompile(`^\d+`)
	matchedNumber, err := strconv.Atoi(countPattern.FindString(event.Text))

	// Ignore messages that don't have numbers.
	if err != nil {
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

	// Retrieve stored info about the last valid number and its sender.
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

	// Ignore numbers that aren't in order.
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

	// Get the current month/year in UTC
	now := time.Now().UTC()
	year := now.Year()
	month := now.Month()

	// Increment the person's monthly count
	redisClient.Incr(fmt.Sprintf("leaderboard:%d-%d:%s", month, year, event.User))
}

func HandleLeaderboardSlashCommand(w http.ResponseWriter, r *http.Request) {
	// Get the current month/year in UTC
	now := time.Now().UTC()
	year := now.Year()
	month := now.Month()

	scan := redisClient.Scan(0, fmt.Sprintf("leaderboard:%d-%d:*", month, year), 10)
	if scan.Err() != nil {
		w.Write([]byte("Something went wrong while loading the leaderboard :cry: Please try again later!"))
		return
	}

	scan_iterator := scan.Iterator()

	type Entry struct {
		Number int
		User   string
	}

	entries := []Entry{}

	for scan_iterator.Next() {
		entry := redisClient.Get(scan_iterator.Val())
		entry_int, err := entry.Int()
		if err != nil {
			return
		}

		if user, ok := parseLeaderboardEntry(scan_iterator.Val()); ok {
			entries = append(entries, Entry{
				Number: entry_int,
				User:   user,
			})
		}
	}

	// Sort entries
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Number > entries[j].Number
	})

	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf(":chart_with_upwards_trend: Counting stats for *%s %d*:", month.String(), year), false, false),
			nil,
			nil,
		),
	}

	for _, v := range entries {
		blocks = append(blocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("<@%s> has counted *%d* this month", v.User, v.Number), false, false), nil, nil))
	}

	resp, _ := json.Marshal(map[string]interface{}{
		"blocks":        blocks,
		"response_type": "in_channel",
	})

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}

func parseLeaderboardEntry(key string) (string, bool) {
	re := regexp.MustCompile(`leaderboard:\d+-\d+:(\w+)`)

	match := re.FindStringSubmatch(key)
	if match == nil {
		return "", false
	}
	return match[1], true
}
