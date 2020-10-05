package pubsub

import (
	"context"

	googlePubSub "cloud.google.com/go/pubsub"

	"github.com/elmagician/pubsub/messages"
)

type (
	// Pubsub provide method to setup and use a pubsub client
	Pubsub interface {
		// Send prepare pubsub to emit a message
		Send(ctx context.Context, msg interface{}) Send
		// Register allow to add Topic or Subscription to Pubsub instance
		Register() Register
		// Listen initialize a lister instance for provided subscription
		Listen() Listener
		// Receive initialize a receiver instance for provided subscription
		Receive() Receiver
		// Clean all known topics and subscriptions
		Clean() error
	}

	// Send interface allow to build messages to be sent through pubsub instance
	Send interface {
		// To indicate topic in witch we would like to send the message
		// if topic is used for the first time, a connection will be created
		// and kept alive regarding this topic
		// Call Clean method to clear all saved topic
		To(topic string) Send
		// WithOption allow to configure locally a send call
		WithOption(opt interface{}) Send
		// Go send message to topics listed in Send instance
		Go() (id string, err error)
		// Destroy indicate to sender that the topics created has to be discarded after
		// publication
		Destroy() Send
	}

	// Register manage known topics and subscriptions
	Register interface {
		// Topic register a new topic using provided publication settings
		// publication settings comes from LINK TO GOOGLE
		Topic(key string, publishOpt *googlePubSub.PublishSettings)
		// Subscription register a subscription topic using provided receive settings
		// receive option comes from LINK TO GOOGLE
		Subscription(key string, receiveOpt *googlePubSub.ReceiveSettings)
	}

	Listener interface {
		// Message initialize a channel to listen to messages of provided Type/Version couple.
		// The provided channel uses the interface type messages.Message but you can
		// safely match it to the provided message type as it is assured that the message emitted in the channel
		// match the type/channel couple
		// newMessage has to be a function witch returns a new Message object. It will be call upon
		// receiving messages to ensure we are using different instances of the messages for each receive messages.
		Message(messageType, messageVersion string, newMessage func() messages.Message) chan messages.Message
		// Unmatched provide a channel to retrieve all messages that could not be matched against provided types/versions.
		Unmatched() chan *messages.Unmatch
		// Error initialize a channel to manage errors
		Error() chan error
		// Start listening
		Start()
		// Stop listening
		Stop()
	}

	Receiver interface {
		// Action apply provided callback method to message matching expected Type && Version
		Action(messageType, messageVersion string, callback func(ctx context.Context, newMessage func() messages.Message))
		// Unmatched apply provided callback to unexpected message Type or Version
		Unmatched(callback func(ctx context.Context, unmatch messages.Unmatch))
		// Error apply a callback method to all received errors
		Error(callback func(ctx context.Context))
		// Start receiving
		Start()
		// Stop receiving
		Stop()
	}
)
