package google

import (
	"context"

	googlePubSub "cloud.google.com/go/pubsub"

	"github.com/elmagician/pubsub"
)

var _ pubsub.SendResults = (*SendResults)(nil)

type SendResults struct {
	results map[string]*googlePubSub.PublishResult
}

func (s SendResults) Results(ctx context.Context) pubsub.Results {
	res := make(pubsub.Results)

	for topic, result := range s.results {
		id, err := result.Get(ctx)
		res[topic] = pubsub.Result{ID: id, Error: err}
	}

	return res
}

func (s SendResults) OnResults(ctx context.Context, callback func(topic string, result pubsub.Result)) {
	for topic, result := range s.results {
		// avoid loop variable overwrite in routing
		result := result
		topic := topic

		go func() {
			id, err := result.Get(ctx)
			callback(topic, pubsub.Result{ID: id, Error: err})
		}()
	}
}
