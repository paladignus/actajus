package nats

import "fmt"

type NATSError struct {
	Op  string
	Err error
	Msg string
}

func (e *NATSError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("nats %s: %s: %v", e.Op, e.Msg, e.Err)
	}
	return fmt.Sprintf("nats %s: %s", e.Op, e.Msg)
}

func (e *NATSError) Unwrap() error {
	return e.Err
}

func ErrConnectionFailed(err error) error {
	return &NATSError{
		Op:  "Connect",
		Err: err,
		Msg: "failed to connect to NATS server",
	}
}

func ErrPublishFailed(err error) error {
	return &NATSError{
		Op:  "Publish",
		Err: err,
		Msg: "failed to publish event",
	}
}

func ErrSubscribeFailed(err error) error {
	return &NATSError{
		Op:  "Subscribe",
		Err: err,
		Msg: "failed to subscribe to events",
	}
}

func ErrStreamCreationFailed(err error) error {
	return &NATSError{
		Op:  "CreateStream",
		Err: err,
		Msg: "failed to create or update stream",
	}
}

func ErrConsumerCreationFailed(err error) error {
	return &NATSError{
		Op:  "CreateConsumer",
		Err: err,
		Msg: "failed to create consumer",
	}
}

func ErrInvalidConfig(msg string) error {
	return &NATSError{
		Op:  "Config",
		Msg: msg,
	}
}

func ErrSerializationFailed(err error) error {
	return &NATSError{
		Op:  "Serialize",
		Err: err,
		Msg: "failed to serialize event to JSON",
	}
}

func ErrDeserializationFailed(err error) error {
	return &NATSError{
		Op:  "Deserialize",
		Err: err,
		Msg: "failed to deserialize event from JSON",
	}
}
