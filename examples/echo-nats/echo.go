package main

import (
	"github.com/nats-io/go-nats"
	"regexp"
	"github.com/rs/zerolog/log"
	"strings"
	"fmt"
	"time"
)

const (
	sub = "slack.event.message"
	pub = "slack.chat.postMessage"
	group = "nats-echo-queue"
)

type NatsClient interface {
	QueueSubscribe(subject, queue string, cb nats.Handler) (*nats.Subscription, error)
	Request(subject string, v interface{}, vPtr interface{}, timeout time.Duration) error
}

func HandleRequest(n NatsClient) {
	n.QueueSubscribe(sub, group, echoHandler(n))
}

func echoHandler(n NatsClient) func (m *slackMsg) {
	return func(m *slackMsg) {
		if m.validateCommand() {

			msg := slackMsg{
				Channel:m.Channel,
				Text: fmt.Sprintf("You said:\n> %s", splitText(m.Text)),
				ThreadTs:m.Ts,
			}
			resp := struct {
				Ok bool `json:"ok"`
				Error string `json:"error"`
			}{}
			if err := n.Request(pub, &msg, &resp, 10 * time.Second); err != nil {
				log.Error().
					Err(err).
					Str("respErr", resp.Error).
					Str("subject", pub).
					Str("message", m.Text).
					Msg("error publishing message to NATS")
			}
		}
	}
}


type slackMsg struct {
	User	string	`json:"user,omitempty"`
	Channel string `json:"channel"`
	Ts		string `json:"ts"`
	ThreadTs string	`json:"thread_ts,omitempty"`
	AsUser	bool	`json:"as_user"`
	Text	string	`json:"text"`
	Type	string	`json:"type"`
}

func (m *slackMsg) validateCommand() bool {
	const rgx = "^nats\\s+echo\\s+.*$"
	ok, _ := regexp.MatchString(rgx, m.Text)
	return ok
}

func splitText(s string) string {
	re := regexp.MustCompile("\\s+")
	sp := re.Split(s, -1)
	return strings.Join(sp[2:], " ")
}