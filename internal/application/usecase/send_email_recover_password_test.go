package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

func TestSendEmailRecoverPassword(t *testing.T) {
	ctx := context.Background()
	logger := &spy.Logger{}
	handler := &spy.Handler{}
	subscriber := &spy.Subscriber{
		EventName:      "user.password_reset_requested",
		SubscribeError: nil,
	}
	t.Run("successful subscription", func(t *testing.T) {
		sut := NewSendEmailRecoverPassword(logger, handler, subscriber)
		logger.On("Info", ctx, "handler registered for event", "event", subscriber.EventName).Once()
		err := sut.Execute(ctx)
		assert.NoError(t, err)
	})
	t.Run("subscription fails", func(t *testing.T) {
		subscriber.SubscribeError = fmt.Errorf("subscription error for testing")
		sut := NewSendEmailRecoverPassword(logger, handler, subscriber)
		logger.On("Error", ctx, "failed to subscribe to event", "event", subscriber.EventName, "error", subscriber.SubscribeError).Once()
		err := sut.Execute(ctx)
		assert.Error(t, err)
	})

	t.Run("logger is called", func(t *testing.T) {
		subscriber.SubscribeError = nil
		sut := NewSendEmailRecoverPassword(logger, handler, subscriber)
		logger.On("Info", ctx, "handler registered for event", "event", subscriber.EventName).Once()
		err := sut.Execute(ctx)
		assert.NoError(t, err)
		assert.True(t, len(logger.Calls) > 0, "expected logger to be called")
	})
}
