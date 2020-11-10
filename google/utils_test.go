package google_test

import (
	"context"

	googlePubSub "cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"google.golang.org/api/option"
	"google.golang.org/grpc"

	"github.com/elmagician/pubsub"
)

func initTestClient(project string, opts ...pstest.ServerReactorOption) (*pstest.Server, *googlePubSub.Client) {
	psTest := pstest.NewServer(opts...)

	conn, err := grpc.Dial(psTest.Addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	cli, err := googlePubSub.NewClient(context.Background(), project, option.WithGRPCConn(conn))
	if err != nil {
		panic(err)
	}

	return psTest, cli
}

type Envelop struct {
	err     error
	message pubsub.Message
	filter  pubsub.MessageFilter
}

func (e *Envelop) ToPubsubMessage() (pubsub.Message, error) {
	if e.err != nil {
		return nil, e.err
	}
	return e.message, nil
}

func (e *Envelop) FromPubsubMessage(msg pubsub.Message) error {
	if e.err != nil {
		e.message = msg
	}
	return e.err
}

func (e *Envelop) Filter() pubsub.MessageFilter {
	return e.filter
}

func (e *Envelop) New() pubsub.Envelop {
	return &Envelop{}
}

type message struct {
	Message *googlePubSub.Message
}

func (m *message) ID() interface{} {
	return m.Message.ID
}

func (m *message) Ack() {
	m.Message.Ack()
}

func (m *message) Nack() {
	m.Message.Nack()
}

func (m *message) Metadata() map[string]string {
	return m.Message.Attributes
}

func (m *message) Data() []byte {
	return m.Message.Data
}

var _ pubsub.Envelop = &Envelop{}

var _ pubsub.Message = &message{}
