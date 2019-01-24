package main

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
	"regexp"
	"strings"
	"time"
)

type NatsClient interface {
	QueueSubscribe(subject, queue string, cb nats.Handler) (*nats.Subscription, error)
	Request(subject string, v interface{}, vPtr interface{}, timeout time.Duration) error
}

func HandleRequest(n NatsClient) {
	n.QueueSubscribe("slack.event.message", "nats-echo-queue", echoHandler(n))
}

func echoHandler(n NatsClient) func(m *slackMsg) {
	return func(m *slackMsg) {
		if ok, _ := regexp.MatchString("^nats echo\\s+.*$", m.Text); !ok {
			return
		}
		msg := slackMsg{
			Channel:  m.Channel,
			Text:     fmt.Sprintf("You said:\n> %s", strings.TrimPrefix(m.Text, "nats echo")),
			ThreadTs: m.Ts,
		}

		var resp slackPostResp
		if err := n.Request("slack.chat.postMessage", &msg, &resp, 10*time.Second); err != nil || resp.Error != "" {
			log.Error().
				Err(err).
				Str("resp", fmt.Sprintf("%+v", resp)).
				Str("subject", "slack.chat.postMessage").
				Str("message", m.Text).
				Msg("error publishing message to NATS")
		}
	}
}

type slackMsg struct {
	Channel  string `json:"channel"`
	Ts       string `json:"ts"`
	ThreadTs string `json:"thread_ts,omitempty"`
	Text     string `json:"text"`
}

type slackPostResp struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}
