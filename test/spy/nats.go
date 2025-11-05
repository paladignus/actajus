// Package spy
package spy

import "context"

type SpyPublisher struct {
	Subject string
	Err     error
}

func (s *SpyPublisher) Publish(ctx context.Context, subject string, msg []byte) error {
	s.Subject = subject
	return s.Err
}
