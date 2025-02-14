package rpc_test

import (
	"bytes"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/synternet/data-layer-sdk/pkg/service"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var _ service.Message = (*Message)(nil)

type Message struct {
	data    []byte
	ch      chan service.Message
	chName  string
	subject string
}

func NewMsg(data []byte, ch chan service.Message, chName, subject string) *Message {
	return &Message{
		data:    data,
		ch:      ch,
		chName:  chName,
		subject: subject,
	}
}

// Ack implements service.Message.
func (m *Message) Ack(opts ...nats.AckOpt) error {
	return nil
}

// AckSync implements service.Message.
func (m *Message) AckSync(opts ...nats.AckOpt) error {
	return nil
}

// Data implements service.Message.
func (m *Message) Data() []byte {
	return m.data
}

// Equal implements service.Message.
func (m *Message) Equal(msg service.Message) bool {
	return bytes.Compare(m.data, msg.Data()) == 0
}

// Header implements service.Message.
func (m *Message) Header() nats.Header {
	return nats.Header{
		"identity": []string{"some", "identity"},
	}
}

// InProgress implements service.Message.
func (m *Message) InProgress(opts ...nats.AckOpt) error {
	panic("unimplemented")
}

// Message implements service.Message.
func (m *Message) Message() *nats.Msg {
	panic("unimplemented")
}

// Metadata implements service.Message.
func (m *Message) Metadata() (*nats.MsgMetadata, error) {
	panic("unimplemented")
}

// Nak implements service.Message.
func (m *Message) Nak(opts ...nats.AckOpt) error {
	return nil
}

// NakWithDelay implements service.Message.
func (m *Message) NakWithDelay(delay time.Duration, opts ...nats.AckOpt) error {
	panic("unimplemented")
}

// QueueName implements service.Message.
func (m *Message) QueueName() string {
	panic("unimplemented")
}

// Reply implements service.Message.
func (m *Message) Reply() string {
	return m.chName
}

// Respond implements service.Message.
func (m *Message) Respond(msg proto.Message) error {
	msgData, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	m.ch <- NewMsg(msgData, nil, "", m.chName)
	return nil
}

// Subject implements service.Message.
func (m *Message) Subject() string {
	return m.subject
}

// Term implements service.Message.
func (m *Message) Term(opts ...nats.AckOpt) error {
	return nil
}
