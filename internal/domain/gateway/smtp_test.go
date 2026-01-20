package gateway

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/test/spy"
)

func TestSMTPInterface(t *testing.T) {
	ctx := context.Background()
	sut := &spy.SMTP{}
	err := sut.SendEmail(ctx, "test@example.com", "Test Subject", "Test Body")
	if err != nil {
		t.Errorf("SendEmail returned error: %v", err)
	}
}
