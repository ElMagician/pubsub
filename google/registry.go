package google

import (
	"context"
	"time"

	googlePubSub "cloud.google.com/go/pubsub"

	"github.com/elmagician/pubsub"
	"github.com/elmagician/pubsub/google/internal"
)

const timeout = 10 * time.Second

var _ pubsub.Registry = (*Registry)(nil)

type Registry struct {
	client        *googlePubSub.Client
	topics        map[string]*internal.Topic
	subscriptions map[string]*internal.Subscription
}

// nolint: dupl
func (l *Registry) AddTopic(key string, publishSettings *googlePubSub.PublishSettings) error {
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
func (l *Registry) AddSubscription(key string, receiveSettings *googlePubSub.ReceiveSettings) error {
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

// nolint: dupl
func (l *Registry) MustAddTopic(key string, publishSettings *googlePubSub.PublishSettings) pubsub.Registry {
	if _, ok := l.topics[key]; ok {
		return nil
	}

	l.topics[key] = &internal.Topic{
		Topic:           l.client.Topic(key),
		PublishSettings: publishSettings,
	}

	return l
}

// nolint: dupl
func (l *Registry) MustAddSubscription(key string, receiveSettings *googlePubSub.ReceiveSettings) pubsub.Registry {
	if _, ok := l.subscriptions[key]; ok {
		return nil
	}

	l.subscriptions[key] = &internal.Subscription{
		Subscription:    l.client.Subscription(key),
		ReceiveSettings: receiveSettings,
	}

	return l
}

func (l *Registry) StopTopics(topics ...string) {
	for _, key := range topics {
		l.StopTopic(key)
	}
}

func (l *Registry) Clear() {
	for key := range l.topics {
		l.StopTopic(key)
	}
	l.topics = make(map[string]*internal.Topic)
	l.subscriptions = make(map[string]*internal.Subscription)
}

func (l *Registry) StopTopic(key string) {
	if topic, ok := l.topics[key]; ok {
		topic.Stop()
	}
}
