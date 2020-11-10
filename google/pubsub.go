package google

import (
	"context"

	googlePubSub "cloud.google.com/go/pubsub"
	"google.golang.org/api/option"

	"github.com/elmagician/pubsub"
	"github.com/elmagician/pubsub/google/internal"
)

var _ pubsub.Pubsub = (*Pubsub)(nil)

// Pubsub implements pubsub.Pubsub interface for GCP.
type Pubsub struct {
	// Client is the gcp instance used to send requests.
	// It is passed as a private parameter to all structures
	// derived from Pubsub
	Client *googlePubSub.Client

	// Config for running instance.
	// It is passed as a private parameter to all structures
	// derived from Pubsub
	Config Config

	registry *Registry
}

// NewPubsub initializes a GCP implementation for pubsub.
func NewPubsub(ctx context.Context, config Config, opts ...option.ClientOption) (pubsub.Pubsub, error) {
	if config.CredentialsPath != "" {
		opts = append(opts, option.WithCredentialsFile(config.CredentialsPath))
	}

	cli, err := googlePubSub.NewClient(ctx, config.ProjectID, opts...)
	if err != nil {
		return nil, err
	}

	return &Pubsub{
			Client: cli, Config: config,
			registry: &Registry{
				client:        cli,
				topics:        make(map[string]*internal.Topic),
				subscriptions: make(map[string]*internal.Subscription),
			},
		},
		nil
}

// Publisher set up an instance to send message to pubsub.
// It use client, config and registry from main Pubsub instance.
func (p Pubsub) Publish() pubsub.Publisher {
	return &Publisher{
		client:    p.Client,
		config:    p.Config,
		registry:  p.registry,
		newTopics: make(map[string]*internal.Topic),
		topics:    make(map[string]*googlePubSub.Topic),
	}
}

func (p Pubsub) Registry() pubsub.Registry {
	return p.registry
}

func (p Pubsub) Listen(subscription string) pubsub.Listener {
	panic("implement me")
}

func (p Pubsub) Receive(subscription string) pubsub.Receiver {
	panic("implement me")
}
