package internal

import (
	"cloud.google.com/go/pubsub"
)

type (
	// Topic provide a structure for pubsub registry
	// it keeps the alive topic Instance linked to a key
	// and its default PublishSettings
	Topic struct {
		*pubsub.Topic
		PublishSettings *pubsub.PublishSettings
	}

	// Subscription provide a structure for pubsub registry
	// it keeps the alive subscription Instance linked to a key
	// and its default ReceiveSettings
	Subscription struct {
		*pubsub.Subscription
		ReceiveSettings *pubsub.ReceiveSettings
	}
)
