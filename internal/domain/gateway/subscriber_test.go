package gateway

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/test/spy"
)

func TestSubscriberInterface(t *testing.T) {
	ctx := context.Background()
	Subscriber := &spy.Subscriber{}
	Handler := &spy.Handler{}
	err := Subscriber.Subscribe(ctx, "test.event", Handler)
	if err != nil {
		t.Errorf("Subscribe returned error: %v", err)
	}
	err = Subscriber.Unsubscribe("test.event")
	if err != nil {
		t.Errorf("Unsubscribe returned error: %v", err)
	}
	err = Subscriber.Close()
	if err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}
