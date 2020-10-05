package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/elmagician/pubsub/messages"
)

var _ messages.Message = &Messages{}

type Messages struct {
	mock.Mock
	ack  func()
	nack func()
}

func NewMessages() *Messages {
	return &Messages{}
}

func (m *Messages) GetMock() *mock.Mock {
	return &m.Mock
}

func (m *Messages) ToPubsubMessage() (metaData map[string]string, data []byte, err error) {
	args := m.Called()
	return args.Get(0).(map[string]string), args.Get(1).([]byte), args.Error(2) // nolint: gomnd
}

func (m *Messages) FromPubsubMessage(metaData map[string]string, data []byte, ack func(), nack func()) error {
	m.ack = ack
	m.nack = nack
	return m.Called(metaData, data, ack, nack).Error(0)
}

func (m *Messages) Ack() {
	m.Called()
	if m.ack != nil {
		m.ack()
	}
}

func (m *Messages) Nack() {
	m.Called()
	if m.nack != nil {
		m.nack()
	}
}

func (m *Messages) Version() string {
	return m.Called().String(0)
}

func (m *Messages) Type() string {
	return m.Called().String(0)
}
