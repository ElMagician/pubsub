package google

import (
	"context"
	"time"

	googlePubSub "cloud.google.com/go/pubsub"

	"github.com/elmagician/pubsub"
	"github.com/elmagician/pubsub/internal"
)

const timeout = 10 * time.Second

var _ pubsub.Registry = (*LocalRegistry)(nil)

type LocalRegistry struct {
	client        *googlePubSub.Client
	topics        map[string]*internal.Topic
	subscriptions map[string]*internal.Subscription
}

// nolint: dupl
func (l *LocalRegistry) AddTopic(key string, publishSettings *googlePubSub.PublishSettings) error {
	if _, ok := l.topics[key]; ok {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	topic := l.client.Topic(key)
	if ok, err := topic.Exists(ctx); !ok || err != nil {
		if err != nil {
			return err
		}
		return pubsub.ErrNotFound
	}

	l.topics[key] = &internal.Topic{
		Topic:           topic,
		PublishSettings: publishSettings,
	}

	return nil
}

// nolint: dupl
func (l *LocalRegistry) AddSubscription(key string, receiveSettings *googlePubSub.ReceiveSettings) error {
	if _, ok := l.subscriptions[key]; ok {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	subscription := l.client.Subscription(key)
	if ok, err := subscription.Exists(ctx); !ok || err != nil {
		if err != nil {
			return err
		}
		return pubsub.ErrNotFound
	}

	l.subscriptions[key] = &internal.Subscription{
		Subscription:    subscription,
		ReceiveSettings: receiveSettings,
	}

	return nil
}

func (l *LocalRegistry) StopTopics(topics ...string) {
	for _, key := range topics {
		l.StopTopic(key)
	}
}

func (l *LocalRegistry) Clear() {
	for key := range l.topics {
		l.StopTopic(key)
	}
	l.topics = make(map[string]*internal.Topic)
	l.subscriptions = make(map[string]*internal.Subscription)
}

func (l *LocalRegistry) StopTopic(key string) {
	if topic, ok := l.topics[key]; ok {
		topic.Stop()
	}
}
