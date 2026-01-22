// Package nats
// This file contains the init/connetc
package nats

import (
	"log/slog"

	"github.com/nats-io/nats.go"
)

func InitJetStream(url string) (nats.JetStreamContext, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		slog.Info("NATS connection fialed", "url", url)
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		slog.Info("NATS connection fialed", "url", url)
		return nil, err
	}

	slog.Info("NATS connected", "url", url)

	return js, err
}
