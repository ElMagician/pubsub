package pubsub

import (
	"context"
	"time"

	googlePubSub "cloud.google.com/go/pubsub"
	"github.com/go-errors/errors"
	"google.golang.org/api/option"

	"github.com/elmagician/pubsub/messages"
)

var (
	_ Pubsub   = &GCP{}
	_ Listener = &GCPListener{}
)

type (
	// Config config object to init a new publisher
	// ProjectID google cloud project
	// JSONConfigPath path to credential json file used for a publisher instance
	// Concurrency represent the number of messages that can be stored in listening channel
	// Timeout set the maximal duration for a pubsub request
	Config struct {
		ProjectID      string        `yaml:"projectId"`
		JSONConfigPath string        `yaml:"configFile"`
		Concurrency    int           `yaml:"concurrency"`
		Timeout        time.Duration `yaml:"timeout"`
	}

	// GCP structure to support Pubsub implementation from Google Cloud Platform
	GCP struct {
		Client       *googlePubSub.Client
		Config       Config
		injectionKey string
	}

	// GCPListener Pubsub listener implementation for GCP
	GCPListener struct {
		client            *googlePubSub.Client
		subscription      *googlePubSub.Subscription
		channelBufferSize int

		chError     chan error
		chUnMatched chan *messages.Unmatch

		ctx       context.Context
		ctxCancel context.CancelFunc
		listeners map[listenerKey]listenerValue
	}

	listenerKey struct {
		messageVersion string
		messageType    string
	}

	listenerValue struct {
		channel chan messages.Message
		message func() messages.Message
	}
)

// New initialize a GCP Pubsub
func New(conf Config, injectionKey string, opts ...option.ClientOption) (pubsub Pubsub, err error) {
	if conf.JSONConfigPath != "" {
		opts = append(opts, option.WithCredentialsFile(conf.JSONConfigPath))
	}
	var cli *googlePubSub.Client
	cli, err = googlePubSub.NewClient(
		context.Background(), conf.ProjectID,
		opts...,
	)
	return &GCP{Client: cli, Config: conf, injectionKey: injectionKey}, err
}

// Send message to topics. Stop on first error
func (gcp *GCP) Send(providedContext context.Context, topics []string, message messages.Message) error {
	ctx, cancel := gcp.ctx(providedContext)
	defer cancel()

	attr, data, err := message.ToPubsubMessage()
	if err != nil {
		return err // error returned from ToPubsubMessage are matchable with Is. No need to wrap them
	}

	messageHandler := &googlePubSub.Message{
		Data:       data,
		Attributes: attr,
	}

	for _, topicKey := range topics {
		resp := gcp.Client.Topic(topicKey).Publish(ctx, messageHandler)
		if _, err := resp.Get(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Listener initialize a pubSub listener for GCP.
// It will fail with unknown subscription if subscription could not be validated from cloud.
func (gcp *GCP) Listener(providedContext context.Context, subscriptionID string) (Listener, error) {
	ctx, cancel := gcp.ctx(providedContext)

	subscription := gcp.Client.Subscription(subscriptionID)
	if exists, err := subscription.Exists(ctx); !exists || err != nil {
		cancel()
		return nil, err
	}

	return &GCPListener{
		subscription:      subscription,
		client:            gcp.Client,
		channelBufferSize: gcp.Config.Concurrency,
		ctx:               ctx,
		ctxCancel:         cancel,
		listeners:         map[listenerKey]listenerValue{},
	}, nil
}

// Implement Godim interface to set injection Key
func (gcp *GCP) Key() string {
	return gcp.injectionKey
}

// Listener to messageType/messageVersion couple using messages.Message transport.
func (listener *GCPListener) Listen(
	messageType, messageVersion string, initMessage func() messages.Message,
) (chMessage chan messages.Message) {
	chMessage = make(chan messages.Message, listener.channelBufferSize)
	listener.listeners[listenerKey{
		messageVersion: messageVersion,
		messageType:    messageType,
	}] = listenerValue{
		channel: chMessage,
		message: initMessage,
	}
	return chMessage
}

// Unmatched allow to receive all message unmatched in a Listener clause.
func (listener *GCPListener) Unmatched() (chUnmatched chan *messages.Unmatch) {
	if listener.chUnMatched == nil {
		listener.chUnMatched = make(chan *messages.Unmatch, listener.channelBufferSize)
	}
	return listener.chUnMatched
}

// Error listen to errors
func (listener *GCPListener) Error() (chError chan error) {
	if listener.chError == nil {
		listener.chError = make(chan error)
	}
	return listener.chError
}

// Start listening to GCP subscription
func (listener *GCPListener) Start() {
	go func() {
		err := listener.subscription.Receive(listener.ctx, func(ctx context.Context, m *googlePubSub.Message) {
			attr := m.Attributes
			version, okVersion := attr["version"]
			msgType, okType := attr["type"]

			// Try to match message to declared transport
			if okVersion && okType {
				listenInfo, ok := listener.listeners[listenerKey{
					messageVersion: version,
					messageType:    msgType,
				}]

				if ok {
					msg := listenInfo.message()
					if err := msg.FromPubsubMessage(attr, m.Data, m.Ack, m.Nack); err == nil {
						listenInfo.channel <- msg
						return
					}
				}
			}

			if listener.chUnMatched != nil {
				listener.chUnMatched <- &messages.Unmatch{
					ID:         m.ID,
					Raw:        m.Data,
					Attributes: attr,
					Ack:        m.Ack,
					Nack:       m.Nack,
				}
			} else {
				m.Ack() // acknowledge unmatched messages by default if they are not listened
			}
		})
		if err != nil && !errors.Is(err, context.Canceled) {
			if listener.chError != nil {
				listener.chError <- err
			}
		}
	}()
}

// Stop listening to GCP subscription
func (listener *GCPListener) Stop() {
	if listener.chError != nil {
		close(listener.chError)
	}
	if listener.chUnMatched != nil {
		close(listener.chUnMatched)
	}
	for _, val := range listener.listeners {
		close(val.channel)
	}
	listener.ctxCancel()
}

func (gcp *GCP) ctx(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, gcp.Config.Timeout)
}
