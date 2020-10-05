package messages

import (
	"regexp"

	v1 "github.com/elmagician/pubsub/messages/v1"
)

var (
	AcceptedSourcesRegex = regexp.MustCompile("^[a-z_-]+$")
)

type (
	// Envelop is an interface to manage pubsub messages relevant data
	// It has to represent the message payload and can include relevant information
	// from message attributes
	Message interface {
		// AsByte convert the envelop to JSON representative byte table
		// Used before emitting message to queue
		ToPubsubMessage() (metaData map[string]string, data []byte, err error)
		// Receive set envelop data from message json payload
		FromPubsubMessage(metaData map[string]string, data []byte, ack func(), nack func()) error
		// Ack acknowledge message
		Ack()
		// Nack refuse to acknowledge message
		Nack()
		// Event version
		Version() string
		// Event type
		Type() string
	}

	Unmatch struct {
		ID         interface{}
		Raw        []byte
		Attributes map[string]string
		Ack        func() `json:"-"`
		Nack       func() `json:"-"`
	}
)

// NewOtherV1 initialize a V1 unspecialized message
func NewOtherV1(source string, marshAbleInterface interface{}) Message {
	return v1.NewOther(source, marshAbleInterface)
}

// NewEmptyOtherV1 initialize an empty V1 unspecialized message
// Used for clearer code when initializing PubSub listeners
func NewEmptyOtherV1() Message {
	return v1.NewEmptyOther()
}
