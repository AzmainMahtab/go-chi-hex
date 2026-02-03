package ports

import (
	"context"

	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
)

type AuditRepository interface {
	Create(ctx context.Context, auditLog domain.Audit) error
}

type AuditPublisher interface {
	Publish(ctx context.Context, auditLog domain.Audit) error
}

type BackgroundWorker interface {
	Start(ctx context.Context) error
	Stop()
}
