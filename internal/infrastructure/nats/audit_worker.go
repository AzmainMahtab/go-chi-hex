package nats

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
	"github.com/AzmainMahtab/go-chi-hex/internal/ports"
	"github.com/nats-io/nats.go"
)

type auditWorker struct {
	js   nats.JetStreamContext
	repo ports.AuditRepository
	wg   sync.WaitGroup
}

func NewAuditWorker(js nats.JetStreamContext, repo ports.AuditRepository) ports.BackgroundWorker {
	return &auditWorker{js: js, repo: repo}
}

func (w *auditWorker) Start(ctx context.Context) error {
	// 'audit-durable' ensures NATS remembers where we left off if we restart
	sub, err := w.js.PullSubscribe("audit.event", "audit-durable")
	if err != nil {
		slog.Error("Error subscbing to event stream", "error", err)
	}

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msgs, err := sub.Fetch(1, nats.MaxWait(1*time.Second))
				if err != nil {
					continue
				}

				msg := msgs[0]
				var e domain.Audit
				json.Unmarshal(msg.Data, &e)

				// SAVE to the Partitioned Postgres table
				if err := w.repo.Create(ctx, e); err != nil {
					msg.Nak() // Try again later
					continue
				}
				msg.Ack() // Success! Delete from NATS
			}
		}
	}()
	return nil
}

func (w *auditWorker) Stop() { w.wg.Wait() }
