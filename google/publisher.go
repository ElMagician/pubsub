package google

import (
	"context"
	"errors"

	googlePubSub "cloud.google.com/go/pubsub"

	"github.com/elmagician/pubsub"
	"github.com/elmagician/pubsub/google/internal"
)

var ErrPublisherDestroyed = errors.New("publisher instance was destroyed")

// Publisher implements pubsub.Publisher interface for GCP.
type Publisher struct {
	client   *googlePubSub.Client
	config   Config
	registry *Registry

	// following value HAS TO BE reset after each call to Send
	destroyed bool
	topics    map[string]*googlePubSub.Topic
	newTopics map[string]*internal.Topic
	nbTopics  int
}

// To adds topic to sent message to.
//
// If Destroy is called, unknown topic will not be saved in registry
// and https://pkg.go.dev/cloud.google.com/go/pubsub#Topic.Stop will be called.
//
// This method apply last registered send configuration for topic. If no configuration where registered, it use default
// configuration
func (p *Publisher) To(topics ...string) pubsub.Publisher {
	defaultConfig := &googlePubSub.PublishSettings{Timeout: p.config.Timeout}

	// add topic to destination
	for _, topicKey := range topics {
		topic, ok := p.registry.topics[topicKey]
		if !ok {
			// Create a new topic with default send configuration if not registered
			topic = &internal.Topic{
				Topic:           p.client.Topic(topicKey),
				PublishSettings: &googlePubSub.PublishSettings{Timeout: p.config.Timeout},
			}

			p.newTopics[topicKey] = topic
		}

		p.topics[topicKey] = topic.Topic

		// Apply last registered configuration for topic
		if topic.PublishSettings != nil {
			topic.Topic.PublishSettings = *topic.PublishSettings
		} else {
			topic.Topic.PublishSettings = *defaultConfig
		}

		p.nbTopics++
	}

	return p
}

// WithOption provide send settings to apply to call.
//
// It will be applied to all topics added to sent process before
// the WithOption call.
func (p *Publisher) WithOption(opt interface{}) pubsub.Publisher {
	panic("implement me")
}

// Send message to topics registered. Any new topic will be saved to the registry if Destroy was not called.
// If Destroy was called, all topics will have their connection stopped and known topics will be kept in registry.
//
// Send will return a single string ID if sending to a single topic, else a list of string.
func (p *Publisher) Send(ctx context.Context, msg pubsub.Envelop) (pubsub.SendResults, error) {
	if p.destroyed {
		return nil, ErrPublisherDestroyed
	}

	res := SendResults{results: make(map[string]*googlePubSub.PublishResult)}

	pubsubMessage, err := msg.ToPubsubMessage() // transform envelop to pubsub.Message
	if err != nil {
		return nil, nil
	}

	// send message to all topics demanded
	for topicKey, topic := range p.topics {
		res.results[topicKey] = topic.Publish(ctx, &googlePubSub.Message{Data: pubsubMessage.Data(), Attributes: pubsubMessage.Metadata()})

		// Add new topics to registry
		if newTopic, ok := p.newTopics[topicKey]; ok {
			p.registry.topics[topicKey] = newTopic
		}
	}

	p.reset()

	return res, nil
}

// TODO clean destroy implementation
func (p *Publisher) Destroy() {
	p.destroyed = true
}

// reset restore default value for variables used for send process.
// It has to be called after each send process.
func (p *Publisher) reset() {
	p.topics = make(map[string]*googlePubSub.Topic)
	p.newTopics = make(map[string]*internal.Topic)
	p.nbTopics = 0
}
