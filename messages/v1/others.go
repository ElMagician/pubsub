package v1

import (
	"encoding/json"
	"time"
)

type Other struct {
	V1
	Payload interface{} `json:"payload"`
}

func NewOther(source string, marshAbleInterface interface{}) *Other {
	return &Other{
		V1: V1{
			EmittedAt: time.Now().UnixNano(),
			Source:    source,
		},
		Payload: marshAbleInterface,
	}
}

// NewEmptyOver return a pointer to empty Other
// It is mainly used to have clearer code when creating PubSub listeners
func NewEmptyOther() *Other {
	return &Other{}
}

func (o *Other) ToPubsubMessage() (metaData map[string]string, data []byte, err error) {
	data, err = json.Marshal(o) // Marshal cannot fail on objects
	return o.V1.ToPubsubMessage(OtherType), data, err
}

func (o *Other) FromPubsubMessage(metaData map[string]string, data []byte, ack func(), nack func()) error {
	if err := json.Unmarshal(data, o); err != nil {
		return err
	}

	return o.V1.FromPubsubMessage(metaData, ack, nack)
}

func (o *Other) Type() string { return OtherType }
