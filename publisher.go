package pubsub

import (
	"context"
)

type (
	// Publisher interface allow to build messages to be sent through pubsub instance.
	Publisher interface {
		// To indicates topic in witch we would like to send the message
		// if topic is used for the first time, a connection will be created
		// and kept alive regarding this topic.
		// Call Clean method to clear all saved topic.
		To(topics ...string) Publisher

		// WithOption allows to configure locally a send call.
		WithOption(opt interface{}) Publisher

		// Send message to topics listed in Send instance. It will returns a SendResults interface
		// witch you can safely discard if you don't need to check that your message
		// was correctly sent.
		Send(ctx context.Context, msg Envelop) (SendResults, error)

		// Destroy has to be called at the end of life of the publisher instance to ensure all messages are correctly
		// sent. Destroy method will only return after ensuring messages were sent or errored then it will
		// destroy connection to pubsub instance definitively.
		// Publisher cannot be used any more after Destroy.
		Destroy()
	}

	// SendResults allows to manage publish results
	SendResults interface {
		// Results recovers send response and return the list of result corresponding Results structure.
		// This is a locking process. Results will await server response before returning.
		Results(ctx context.Context) Results

		// OnResults will apply callback function when server respond.
		OnResults(ctx context.Context, allback func(topic string, result Result))
	}

	Result struct {
		ID    string
		Error error
	}

	Results map[string]Result
)
