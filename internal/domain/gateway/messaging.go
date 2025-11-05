// Package gateway
package gateway

import "context"

// Publisher publica mensagens em um "subject" (ex: "user.created")
type Publisher interface {
	Publish(ctx context.Context, subject string, msg []byte) error
}
