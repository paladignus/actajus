// Package messaging
package messaging

import "context"

type Publisher interface {
	Publish(ctx context.Context, subject string, idMessage string, payload []byte) error
}
