package gateway

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/test/spy"
)

// MockHandler is a mock implementation of the repository.IHandler interface for testing

// func (m *MockHandler) Handle(ctx context.Context, event event.IEvent) error {
// 	return nil
// }
//
// func (m *MockHandler) CanHandle(event event.IEvent) bool {
// 	return true
// }
//
// // MockSubscriber is a mock implementation of the Subscriber interface for testing
// type MockSubscriber struct{}
//
// func (m *MockSubscriber) Subscribe(ctx context.Context, eventName string, handler repository.IHandler) error {
// 	return nil
// }
//
// func (m *MockSubscriber) Unsubscribe(eventName string) error {
// 	return nil
// }
//
// func (m *MockSubscriber) Close() error {
// 	return nil
// }

// TestSubscriberInterface ensures that our interface is correctly defined
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
