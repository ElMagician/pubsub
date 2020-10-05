package v1

import (
	"encoding/gob"
	stdError "errors"
	"strconv"
)

const (
	VersionKey = "v1"
	AuditType  = "audit"
	EventType  = "event"
	OtherType  = "other"
)

// Register type for gob serializer
func init() {
	gob.Register(V1{})
	gob.Register(Other{})
}

type V1 struct {
	EmittedAt int64  `json:"-"`
	Source    string `json:"-"`
	ack       func()
	nack      func()
}

func (v1 *V1) ToPubsubMessage(msgType string) (metaData map[string]string) {
	metaData = map[string]string{
		"version":   VersionKey,
		"type":      msgType,
		"source":    v1.Source,
		"emittedAt": strconv.FormatInt(v1.EmittedAt, 10),
	}
	return metaData
}

func (v1 *V1) FromPubsubMessage(metaData map[string]string, ack func(), nack func()) error {
	v1.ack = ack
	v1.nack = nack
	src := "unknow"

	if source, ok := metaData["source"]; ok {
		src = source
	}

	emit, ok := metaData["emittedAt"]
	if !ok {
		return stdError.New("emittedAt missing") // nolint:goerr113
	}
	emitInt, err := strconv.ParseInt(emit, 10, 64)
	if err != nil {
		return err // nolint:goerr113
	}

	v1.EmittedAt = emitInt
	v1.Source = src

	return nil
}

func (v1 *V1) Ack()         { v1.ack() }
func (v1 *V1) Nack()        { v1.nack() }
func (*V1) Version() string { return VersionKey }
