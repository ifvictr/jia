package jia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var (
	jiaConfig   *Config
	redisClient *redis.Client
	slackClient *slack.Client
)

func StartServer(config *Config) {
	jiaConfig = config
	// Set up Redis connection
	options, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		panic(err)
	}
	redisClient = redis.NewClient(options)

	// Initialize Slack app
	slackClient = slack.New(config.BotToken)

	// Start receiving events
	http.HandleFunc("/slack/events", handleSlackEvents)
	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", config.Port), nil)
}

func handleSlackEvents(w http.ResponseWriter, r *http.Request) {
	// Verify the payload was sent by Slack.
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	apiEvent, err := slackevents.ParseEvent(json.RawMessage(body),
		slackevents.OptionVerifyToken(
			&slackevents.TokenComparator{VerificationToken: jiaConfig.VerificationToken}))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Handle the event that came through
	switch apiEvent.Type {
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
		break
	case slackevents.CallbackEvent:
		HandleInnerEvent(slackClient, &apiEvent.InnerEvent)
		break
	}
}
