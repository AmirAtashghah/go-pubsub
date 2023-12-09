package message

import (
	"context"
)

// Publisher is the emitting part of a Pub/Sub.
type Publisher interface {
	Publish(topic string, messages ...*Message) error
	// Close should flush unsent messages, if publisher is async.
	Close() error
}

// Subscriber is the consuming part of the Pub/Sub.
type Subscriber interface {
	Subscribe(ctx context.Context, topic string) (<-chan *Message, error)
	// Close closes all subscriptions with their output channels and flush offsets etc. when needed.
	Close() error
}
