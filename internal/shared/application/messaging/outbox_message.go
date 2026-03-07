// Package messaging
package messaging

import "time"

type OutboxMessage struct {
	ID          string
	Subject     string
	Payload     []byte
	OccurredAt  time.Time
	PublishedAt *time.Time
}
