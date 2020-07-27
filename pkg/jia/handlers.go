package jia

import (
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
}
