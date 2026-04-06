package nats

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNATSError_Error(t *testing.T) {
	t.Run("with underlying error", func(t *testing.T) {
		underlyingErr := errors.New("underlying error")
		err := &NATSError{
			Op:  "Test",
			Err: underlyingErr,
			Msg: "test message",
		}

		expected := "nats Test: test message: underlying error"
		actual := err.Error()

		assert.Equal(t, expected, actual)
	})

	t.Run("without underlying error", func(t *testing.T) {
		err := &NATSError{
			Op:  "Test",
			Err: nil,
			Msg: "test message",
		}

		expected := "nats Test: test message"
		actual := err.Error()

		assert.Equal(t, expected, actual)
	})
}

func TestNATSError_Unwrap(t *testing.T) {
	t.Run("with underlying error", func(t *testing.T) {
		underlyingErr := errors.New("underlying error")
		err := &NATSError{
			Op:  "Test",
			Err: underlyingErr,
			Msg: "test message",
		}

		unwrapped := err.Unwrap()

		assert.Equal(t, underlyingErr, unwrapped)
	})

	t.Run("without underlying error", func(t *testing.T) {
		err := &NATSError{
			Op:  "Test",
			Err: nil,
			Msg: "test message",
		}

		unwrapped := err.Unwrap()

		assert.Nil(t, unwrapped)
	})
}

func TestErrConnectionFailed(t *testing.T) {
	underlyingErr := errors.New("connection failed")
	err := ErrConnectionFailed(underlyingErr)

	natsErr, ok := err.(*NATSError)
	assert.True(t, ok)
	assert.Equal(t, "Connect", natsErr.Op)
	assert.Equal(t, "failed to connect to NATS server", natsErr.Msg)
	assert.Equal(t, underlyingErr, natsErr.Err)
}

func TestErrPublishFailed(t *testing.T) {
	underlyingErr := errors.New("publish failed")
	err := ErrPublishFailed(underlyingErr)

	natsErr, ok := err.(*NATSError)
	assert.True(t, ok)
	assert.Equal(t, "Publish", natsErr.Op)
	assert.Equal(t, "failed to publish event", natsErr.Msg)
	assert.Equal(t, underlyingErr, natsErr.Err)
}

func TestErrSubscribeFailed(t *testing.T) {
	underlyingErr := errors.New("subscribe failed")
	err := ErrSubscribeFailed(underlyingErr)

	natsErr, ok := err.(*NATSError)
	assert.True(t, ok)
	assert.Equal(t, "Subscribe", natsErr.Op)
	assert.Equal(t, "failed to subscribe to events", natsErr.Msg)
	assert.Equal(t, underlyingErr, natsErr.Err)
}

func TestErrStreamCreationFailed(t *testing.T) {
	underlyingErr := errors.New("stream creation failed")
	err := ErrStreamCreationFailed(underlyingErr)

	natsErr, ok := err.(*NATSError)
	assert.True(t, ok)
	assert.Equal(t, "CreateStream", natsErr.Op)
	assert.Equal(t, "failed to create or update stream", natsErr.Msg)
	assert.Equal(t, underlyingErr, natsErr.Err)
}

func TestErrConsumerCreationFailed(t *testing.T) {
	underlyingErr := errors.New("consumer creation failed")
	err := ErrConsumerCreationFailed(underlyingErr)

	natsErr, ok := err.(*NATSError)
	assert.True(t, ok)
	assert.Equal(t, "CreateConsumer", natsErr.Op)
	assert.Equal(t, "failed to create consumer", natsErr.Msg)
	assert.Equal(t, underlyingErr, natsErr.Err)
}

func TestErrInvalidConfig(t *testing.T) {
	err := ErrInvalidConfig("invalid configuration")

	natsErr, ok := err.(*NATSError)
	assert.True(t, ok)
	assert.Equal(t, "Config", natsErr.Op)
	assert.Equal(t, "invalid configuration", natsErr.Msg)
	assert.Nil(t, natsErr.Err)
}

func TestErrSerializationFailed(t *testing.T) {
	underlyingErr := errors.New("serialization failed")
	err := ErrSerializationFailed(underlyingErr)

	natsErr, ok := err.(*NATSError)
	assert.True(t, ok)
	assert.Equal(t, "Serialize", natsErr.Op)
	assert.Equal(t, "failed to serialize event to JSON", natsErr.Msg)
	assert.Equal(t, underlyingErr, natsErr.Err)
}

func TestErrDeserializationFailed(t *testing.T) {
	underlyingErr := errors.New("deserialization failed")
	err := ErrDeserializationFailed(underlyingErr)

	natsErr, ok := err.(*NATSError)
	assert.True(t, ok)
	assert.Equal(t, "Deserialize", natsErr.Op)
	assert.Equal(t, "failed to deserialize event from JSON", natsErr.Msg)
	assert.Equal(t, underlyingErr, natsErr.Err)
}