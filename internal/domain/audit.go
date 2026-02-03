// Package domain
// this package contains the domain for audit
package domain

type Audit struct {
	UUID      string `db:"uuid"`
	EventType string `db:"event_type"`
	ActorID   string `db:"actor_id"`
	Payload   string `db:"payload"`
}
