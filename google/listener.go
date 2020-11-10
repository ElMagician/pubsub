package google

import (
	"context"

	"github.com/elmagician/pubsub"
)

var _ pubsub.Listener = (*Listener)(nil)

type Listener struct{}

func (l Listener) OnMessage(envelop pubsub.Envelop, newMessage func() pubsub.Message) chan pubsub.Message {
	panic("implement me")
}

func (l Listener) OnUnmatched() chan pubsub.Message {
	panic("implement me")
}

func (l Listener) OnError() chan error {
	panic("implement me")
}

func (l Listener) Listen(ctx context.Context) {
	panic("implement me")
}

func (l Listener) Stop() {
	panic("implement me")
}
