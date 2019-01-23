package main

import (
	"os"
	"github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
)

func main() {
	natsURL, ok := os.LookupEnv("NATS_URL")
	if !ok {
		natsURL = nats.DefaultURL
	}

	n := newNatsConn(natsURL)
	defer n.Close()
	go HandleRequest(n)

	select {}
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

