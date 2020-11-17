package pubsub

import (
	"context"
)

type (
	// Listener provides method to setup listening process on subscription. Listened messages will
	// be transformed to a Message interface and sent through channel.
	Listener interface {
		// OnMessage initializes a channel to listen to messages of provided Type/Version couple.
		// The provided channel uses the interface type messages.Message but you can
		// safely match it to the provided message type as it is assured that the message emitted in the channel
		// match the type/channel couple
		// newMessage has to be a function witch returns a new Message object. It will be call upon
		// receiving messages to ensure we are using different instances of the messages for each receive messages.
		OnMessage(envelop Envelop, newMessage func() Message) chan Message

		// OnUnmatched provides a channel to retrieve all messages that could not be matched against provided types/versions.
		OnUnmatched() chan Message

		// OnError initializes a channel to manage errors.
		OnError() chan error

		// Listen starts listening process in background.
		Listen(ctx context.Context)

		// Stop listening.
		Stop()
	}

	// Receiver provides method to setup reception process on subscription. Received messages will
	// be transformed to a Message then process through provided processes.
	Receiver interface {
		// OnMessage applies provided callback method to message matching expected Type && Version.
		OnMessage(envelop Envelop, callback MessageCallback)

		// OnUnmatched applies provided callback to unexpected message Type or Version.
		OnUnmatched(callback MessageCallback)

		// OnError applies a callback method to all received errors.
		// Message will be nil if errors occurs without receiving message from pubsub
		// client.
		OnError(callback func(err error, msg Message))

		// Start receiving as separated process.
		// Errors are managed through Error callback.
		Start(ctx context.Context, reloadOnFailure bool)

		// Receive messages in current process.
		// Process will stop at the first error received and return it.
		// No errors are returned when Stop is used.
		Receive(ctx context.Context) error

		// Stop receiving.
		Stop()
	}
)
