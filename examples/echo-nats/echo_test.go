package main

import (
	"github.com/nats-io/go-nats"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type NatsStub struct {
	Sub   string
	Resp  interface{}
	Reply interface{}
}

func (n *NatsStub) Request(subject string, v interface{}, vPtr interface{}, timeout time.Duration) error {
	n.Sub = subject
	n.Resp = v
	return nil
}

func (n *NatsStub) QueueSubscribe(subject, queue string, cb nats.Handler) (*nats.Subscription, error) {
	return nil, nil
}

func TestHandleRequestSuccess(t *testing.T) {
	// given
	m := slackMsg{
		Text: "nats echo blahblah",
	}
	e := "You said:\n> blahblah"
	eSub := "slack.chat.postMessage"
	var n = NatsStub{}
	// when
	echoHandler(&n)(&m)
	a := n.Resp.(*slackMsg).Text
	aSub := n.Sub
	// then
	assert.Equal(t, e, a)
	assert.Equal(t, eSub, aSub)
}
