// Package nats
// EventPublisher implementation using NATS JetStream
package nats

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
	"github.com/nats-io/nats.go"
)

type NatsEventPublisher struct {
	js nats.JetStreamContext
}

func NewNatsEventPublisher(js nats.JetStreamContext) *NatsEventPublisher {
	return &NatsEventPublisher{js: js}
}

func (n *NatsEventPublisher) Publish(ctx context.Context, audit domain.Audit) error {
	data, err := json.Marshal(audit)
	if err != nil {
		slog.Error("Failed to marshal audit event", "error", err)
		return err
	}

	_, err = n.js.Publish("audit.event", data, nats.Context(ctx))
	slog.Info("Published audit event", "event", err)
	return err
}
