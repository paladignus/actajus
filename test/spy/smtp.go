// Package spy
package spy

import "context"

type SpySMTP struct {
	Email string
	Err   error
	// To     []string
	// From   string
	// Body   string
	// Subject string
}

func (s *SpySMTP) SendEmail(ctx context.Context, to string, resetURL string) error {
	s.Email = to
	return s.Err
}
