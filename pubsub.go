package pubsub

import (
	googlePubSub "cloud.google.com/go/pubsub"
)

type (
	// Pubsub provides method to setup and use a pubsub client.
	Pubsub interface {
		// Publish prepare pubsub to emit a message.
		Publish() Publisher

		// Registry allow to add Topic or Subscription to Pubsub instance.
		Registry() Registry

		// Listen initialize a lister instance for provided subscription.
		Listen(subscription string) Listener

		// Receive initialize a receiver instance for provided subscription.
		Receive(subscription string) Receiver
	}

	// Registry manages known topics and subscriptions.
	Registry interface {
		// AddTopic registers a new topic using provided publication settings.
		AddTopic(key string, publishSettings *googlePubSub.PublishSettings) error

		// MustAddTopic registers a new topic using provided publication settings or panic.
		MustAddTopic(key string, publishSettings *googlePubSub.PublishSettings) Registry

		// AddSubscription registers a subscription topic using provided receive settings.
		AddSubscription(key string, receiveSettings *googlePubSub.ReceiveSettings) error

		// MustAddSubscription registers a subscription topic using provided receive settings or panic.
		MustAddSubscription(key string, receiveSettings *googlePubSub.ReceiveSettings) Registry

		// StopTopics has to be called to kill connection to topic instance. Passing no arguments
		// will stop all known topics.
		StopTopics(topics ...string)

		// Clear registry of all known Topics && Subscriptions. Clear will stop all topics removed.
		Clear()
	}
)
