package pubsub

import (
	"context"

	googlePubSub "cloud.google.com/go/pubsub"
)

type (
	// Pubsub provide method to setup and use a pubsub client
	Pubsub interface {
		// Publisher prepare pubsub to emit a message
		Publisher(ctx context.Context, msg interface{}) Publisher

		// Registry allow to add Topic or Subscription to Pubsub instance
		Registry() Registry

		// Listen initialize a lister instance for provided subscription
		Listen() Listener

		// Receive initialize a receiver instance for provided subscription
		Receive() Receiver

		// Clean all known topics and subscriptions
		Clean() error
	}

	// Registry manage known topics and subscriptions
	Registry interface {
		// AddTopic register a new topic using provided publication settings
		// publication settings comes from LINK TO GOOGLE
		// It returns an error if topic does not exists or pubsub client call failed
		AddTopic(key string, publishSettings *googlePubSub.PublishSettings) error

		// AddSubscription register a subscription topic using provided receive settings
		// receive option comes from LINK TO GOOGLE
		// It returns an error if subscription does not exists or pubsub client call failed
		AddSubscription(key string, receiveSettings *googlePubSub.ReceiveSettings) error

		// StopTopics has to be called to kill connection to topic instance
		StopTopics(topics ...string)

		// Clear registry of all known Topics && Subscriptions. Clear will stop all topics removed
		Clear()
	}

	// Publisher interface allow to build messages to be sent through pubsub instance
	Publisher interface {
		// To indicate topic in witch we would like to send the message
		// if topic is used for the first time, a connection will be created
		// and kept alive regarding this topic
		// Call Clean method to clear all saved topic
		To(topic string) Publisher

		// WithOption allow to configure locally a send call
		WithOption(opt interface{}) Publisher

		// Go send message to topics listed in Send instance
		Send() (id string, err error)

		// Destroy indicate to sender that the topics created has to be discarded after
		// publication
		Destroy() Publisher
	}

	// Listener provide method to setup listening process on subscription. Listened messages will
	// be transformed to a Message interface and sent through channel
	Listener interface {
		// Message initialize a channel to listen to messages of provided Type/Version couple.
		// The provided channel uses the interface type messages.Message but you can
		// safely match it to the provided message type as it is assured that the message emitted in the channel
		// match the type/channel couple
		// newMessage has to be a function witch returns a new Message object. It will be call upon
		// receiving messages to ensure we are using different instances of the messages for each receive messages.
		Message(messageType, messageVersion string, newMessage func() Message) chan Message

		// Unmatched provide a channel to retrieve all messages that could not be matched against provided types/versions.
		Unmatched() chan *Unmatch

		// Error initialize a channel to manage errors
		Error() chan error

		// Start listening
		Start()

		// Stop listening
		Stop()
	}

	// Receiver provide method to setup reception process on subscription. Received messages will
	// be transformed to a Message then process through provided processes.
	Receiver interface {
		// Action apply provided callback method to message matching expected Type && Version
		Action(messageType, messageVersion string, callback func(ctx context.Context, newMessage func() Message))

		// Unmatched apply provided callback to unexpected message Type or Version
		Unmatched(callback func(ctx context.Context, unmatch Unmatch))

		// Error apply a callback method to all received errors
		Error(callback func(ctx context.Context))

		// Start receiving
		Start()

		// Stop receiving
		Stop()
	}
)
