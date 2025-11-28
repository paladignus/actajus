// Package spy
package spy

import "context"

type SMTP struct {
	Email string
	Err   error
	// To     []string
	// From   string
	// Body   string
	// Subject string
}

func (s *SMTP) SendEmail(ctx context.Context, to, subject, body string) error {
	s.Email = to
	return s.Err
}
