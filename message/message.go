package message

import (
	"context"

	"github.com/google/uuid"
)

type payload []byte

type Message struct {
	Id       uuid.UUID
	Payload  payload
	Metadata Metadata

	ctx context.Context
}

// NewMessage creates a new Message with given uuid and payload.
func NewMessage(id uuid.UUID, payload payload) *Message {
	return &Message{
		Id:       id,
		Metadata: make(map[string]string),
		Payload:  payload,
	}
}

// SetContext sets provided context to the message.
func (m *Message) SetContext(ctx context.Context) {
	m.ctx = ctx
}

// Copy copies all message without Acks/Nacks.
// The context is not propagated to the copy.
func (m *Message) Copy() *Message {
	msg := NewMessage(m.Id, m.Payload)
	for k, v := range m.Metadata {
		msg.Metadata.Set(k, v)
	}
	return msg
}
