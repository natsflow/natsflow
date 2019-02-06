package main

import (
	"github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
	"os"
	"runtime"
	"time"
	"regexp"
	"fmt"
	"strings"
)


func main() {
	natsURL, ok := os.LookupEnv("NATS_URL")
	if !ok {
		natsURL = nats.DefaultURL
	}

	n := newNatsConn(natsURL)
	n.QueueSubscribe("slack.event.message", "nats-echo-queue", func(m *slackMsg) {
		if ok, _ := regexp.MatchString("^nats echo\\s+.*$", m.Text); !ok {
			return
		}
		msg := slackMsg{
			Channel:  m.Channel,
			Text:     fmt.Sprintf("You said:\n> %s", strings.TrimPrefix(m.Text, "nats echo ")),
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
	})

	runtime.Goexit()
}

func newNatsConn(host string) *nats.EncodedConn {
	nc, err := nats.Connect(host)
	if err != nil {
		log.Fatal().Err(err).Str("host", host).Msg("could not connect to NATS")
	}

	log.Info().Str("host", host).Msg("connected to NATS")

	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create encoded NATS connection")
	}
	return c
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
