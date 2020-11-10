package google

import (
	"context"

	"github.com/elmagician/pubsub"
)

var _ pubsub.Receiver = (*Receiver)(nil)

type Receiver struct{}

func (r Receiver) OnMessage(envelop pubsub.Envelop, callback pubsub.MessageCallback) {
	panic("implement me")
}

func (r Receiver) OnUnmatched(callback pubsub.MessageCallback) {
	panic("implement me")
}

func (r Receiver) OnError(callback func(ctx context.Context)) {
	panic("implement me")
}

func (r Receiver) Start(ctx context.Context) {
	panic("implement me")
}

func (r Receiver) Receive(ctx context.Context) error {
	panic("implement me")
}

func (r Receiver) Stop() {
	panic("implement me")
}
