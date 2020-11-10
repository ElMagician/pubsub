package google

import (
	googlePubSub "cloud.google.com/go/pubsub"

	"github.com/elmagician/pubsub"
)

var _ pubsub.Message = &message{}

// message provide a private structure to match pubsub.Message (elMagician) interface from pubsub.Message (Google) structure
type message struct {
	Message *googlePubSub.Message
}

func (m *message) ID() interface{} {
	return m.Message.ID
}

func (m *message) Ack() {
	m.Message.Ack()
}

func (m *message) Nack() {
	m.Message.Nack()
}

func (m *message) Metadata() map[string]string {
	return m.Message.Attributes
}

func (m *message) Data() []byte {
	return m.Message.Data
}
