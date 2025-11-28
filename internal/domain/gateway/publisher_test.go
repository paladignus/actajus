package gateway

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/internal/domain/event"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

func TestPublisherInterface(t *testing.T) {
	ctx := context.Background()
	mockPublisher := &spy.SpyPublisher{}
	mockEvent := spy.MockEvent{}
	err := mockPublisher.Publish(ctx, mockEvent)
	assert.NoError(t, err)
	err = mockPublisher.PublishBatch(ctx, []event.IEvent{mockEvent})
	assert.NoError(t, err)
	err = mockPublisher.Close()
	assert.NoError(t, err)
}

