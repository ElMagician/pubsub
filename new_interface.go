package pubsub

import (
	"context"

	"github.com/elmagician/pubsub/messages"
)

type (
	PubSub2 interface {
		// Send msg to a list of topics. If topic is used for the first time,
		// it will be created and registered until Clean method is called.
		Send(ctx context.Context, msg interface{}, topics ...string)
		// SendOnce send a msg to a list of topic witch will be created then clean on this call.
		SendOnce(ctx context.Context, msg interface{}, topics ...string)
		// Listener to a subscription
		Listener(ctx context.Context, subscription string) (Listener2, error)
		// Registry describe an action registry for provided subscription
		Registry(ctx context.Context, subscription string) (ActionRegistry, error)
		// Clean all known topics and subscriptions
		Clean() error
	}

	Listener2 interface {
		// Listener initialize a channel to listen to messages of provided Type/Version couple.
		// The provided channel uses the interface type messages.Message but you can
		// safely match it to the provided message type as it is assured that the message emitted in the channel
		// match the type/channel couple
		// newMessage has to be a function witch returns a new Message object. It will be call upon
		// receiving messages to ensure we are using different instances of the messages for each receive messages.
		Listen(messageType, messageVersion string, newMessage func() messages.Message) chan messages.Message
		// Unmatched provide a channel to retrieve all messages that could not be matched against provided types/versions.
		Unmatched() chan *messages.Unmatch
		// Error initialize a channel to manage errors
		Error() chan error
		// Start listening
		Start()
		// Stop listening
		Stop()
	}

	ActionRegistry interface {
		// Register set callback method for specific usage
		Register() Register2
		// Use registered action
		Use()
		// Discard registered action
		Discard()
	}

	Register2 interface {
		// Action apply provided callback method to message matching expected Type && Version
		Action(messageType, messageVersion string, callback func(ctx context.Context, newMessage func() messages.Message))
		// Unmatched apply provided callback to unexpected message Type or Version
		Unmatched(callback func(ctx context.Context, unmatch messages.Unmatch))
		// Error apply a callback method to all received errors
		Error(callback func(ctx context.Context))
	}
)
