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
	sut := &spy.Publisher{}
	spyEvent := spy.Event{}
	err := sut.Publish(ctx, spyEvent)
	assert.NoError(t, err)
	err = sut.PublishBatch(ctx, []event.IEvent{spyEvent})
	assert.NoError(t, err)
	err = sut.Close()
	assert.NoError(t, err)
}
