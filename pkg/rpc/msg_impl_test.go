package rpc_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/synternet/data-layer-sdk/pkg/service"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var _ service.Message = (*Message)(nil)

type Message struct {
	t            *testing.T
	data         []byte
	ch           chan service.Message
	replySubject string
	subject      string
}

func NewMsg(t *testing.T, data []byte, ch chan service.Message, replySubject, subject string) *Message {
	return &Message{
		t:            t,
		data:         data,
		ch:           ch,
		replySubject: replySubject,
		subject:      subject,
	}
}

// Ack implements service.Message.
func (m *Message) Ack(opts ...nats.AckOpt) error {
	m.t.Log("msg Ack")
	return nil
}

// AckSync implements service.Message.
func (m *Message) AckSync(opts ...nats.AckOpt) error {
	m.t.Log("msg AckSync")
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
	m.t.Log("msg Nak")
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
	return m.replySubject
}

// Respond implements service.Message.
func (m *Message) Respond(msg proto.Message) error {
	msgData, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	m.t.Logf("message Respond: reply=%s subj=%s msg=%s ch=%v", m.replySubject, m.subject, reflect.TypeOf(msg).Name(), m.ch)
	m.ch <- NewMsg(m.t, msgData, nil, "", m.replySubject)
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
