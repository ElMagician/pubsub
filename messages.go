package pubsub

type (
	// Message is an interface to manage pubsub messages relevant data
	// It has to represent the message payload and can include relevant information
	// from message attributes. It is used to abstract messages type and organisation
	// allowing any struct to be converted to a pubsub message while masking
	// some pubsub logic.
	Message interface {
		// ID of the message in the service provider
		ID() interface{}
		// Data payload
		Data() []byte
		// Metadata are all the tags witch can identify the message and group it with others but are not
		// relevant as information
		Metadata() map[string]string
		// Ack acknowledge message
		Ack()
		// Nack refuse to acknowledge message
		Nack()
	}

	// Envelop is an interface to link a golang object representing your
	// relevant data to a Message interface. It act as a DTO and
	// use Filter method to get relevant filtering data.
	Envelop interface {
		// ToPubsubMessage convert the envelop to JSON representative byte table
		// Used before emitting message to queue
		ToPubsubMessage() (Message, error)
		// FromPubsubMessage set envelop data from message json payload
		FromPubsubMessage(msg Message, ack func(), nack func()) error
		// Filter return a logical map to filter value. The key as to match the expected pubsub metadata key
		// while the value represent the expected value to filter
		Filter() map[string]string
	}

	// Unmatch represent messages not matching an expected couple of Version/Type or the Message interface
	// It provides the id of the message and full data as well as methods to Acknowledge or not the message
	Unmatch struct {
		ID         interface{}
		Raw        []byte
		Attributes map[string]string
		Ack        func() `json:"-"`
		Nack       func() `json:"-"`
	}
)
