package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/elmagician/pubsub"
	"github.com/elmagician/pubsub/messages"
)

var (
	_ pubsub.Listener = &Listener{}
	_ pubsub.Pubsub   = &PubSub{}
)

type (
	PubSub struct {
		mock.Mock
	}

	Listener struct {
		mock.Mock
	}
)

func NewPubsub() *PubSub {
	return &PubSub{}
}

func NewListener() *Listener {
	return &Listener{}
}

func (m *PubSub) GetMock() *mock.Mock {
	return &m.Mock
}

func (m *Listener) GetMock() *mock.Mock {
	return &m.Mock
}

func (m *PubSub) Send(ctx context.Context, topic []string, message messages.Message) error {
	args := m.Called(ctx, topic, message)
	return args.Error(0)
}

func (m *PubSub) Listener(ctx context.Context, subscription string) (pubsub.Listener, error) {
	args := m.Called(ctx, subscription)
	return args.Get(0).(pubsub.Listener), args.Error(1)
}

func (m *Listener) Listen(
	messageType, messageVersion string, initMessage func() messages.Message,
) chan messages.Message {
	return m.Called(messageType, messageVersion, initMessage()).Get(0).(chan messages.Message)
}

func (m *Listener) Unmatched() chan *messages.Unmatch {
	return m.Called().Get(0).(chan *messages.Unmatch)
}

func (m *Listener) Error() chan error {
	return m.Called().Get(0).(chan error)
}

func (m *Listener) Start() {
	m.Called()
}

func (m *Listener) Stop() {
	m.Called()
}
