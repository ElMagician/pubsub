package pubsub

import (
	"context"
)

type (
	// Message is an interface to manage pubsub messages relevant data
	// It has to represent the message payload and can include relevant information
	// from message attributes. It is used to abstract messages type and organisation
	// allowing any struct to be converted to a pubsub message while masking
	// some pubsub logic.
	Message interface {
		// ID of the message in the service provider.
		ID() interface{}

		// Data payload.
		Data() []byte

		// Metadata are all the tags witch can identify the message and group it with others but are not
		// relevant as information.
		Metadata() map[string]string

		// Ack acknowledges message.
		Ack()

		// Nack refuses to acknowledge message.
		Nack()
	}

	// Envelop is an interface to link a golang object representing your
	// relevant data to a Message interface. It act as a DTO and
	// use Filter method to get relevant filtering data.
	Envelop interface {
		// ToPubsubMessage converts the envelop to JSON representative byte table.
		// Used before emitting message to queue.
		ToPubsubMessage() (Message, error)

		// FromPubsubMessage set envelop data from message json payload.
		FromPubsubMessage(msg Message) error

		// Filter returns a logical map to filter value. The key as to match the expected pubsub metadata key
		// while the value represent the expected value to filter.
		Filter() MessageFilter

		// New generates a new empty envelop.
		// Used form message reception.
		New() Envelop
	}

	// MessageCallback enforces callback function type on reception.
	MessageCallback func(ctx context.Context, msg Message)

	// MessageFilter represents filtering data for message reception.
	MessageFilter map[string]string
)
