package service

import (
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/synternet/data-layer-sdk/pkg/options"
	"google.golang.org/protobuf/proto"
)

type Message interface {
	Ack(opts ...nats.AckOpt) error
	AckSync(opts ...nats.AckOpt) error
	Equal(msg Message) bool
	InProgress(opts ...nats.AckOpt) error
	Metadata() (*nats.MsgMetadata, error)
	Nak(opts ...nats.AckOpt) error
	NakWithDelay(delay time.Duration, opts ...nats.AckOpt) error
	Respond(proto.Message) error
	Term(opts ...nats.AckOpt) error

	Message() *nats.Msg
	Subject() string
	Reply() string
	Data() []byte
	Header() nats.Header
	QueueName() string
}

type (
	MessageHandler func(msg Message)
	ServiceHandler func(msg Message) (proto.Message, error)
)

type natsMessage struct {
	*nats.Msg
	codec        options.Codec
	msgCounter   *atomic.Uint64
	bytesCounter *atomic.Uint64
	make         func([]byte, string, string) (*nats.Msg, error)
}

func wrapMessage(codec options.Codec, msgCounter, bytesCounter *atomic.Uint64, maker func([]byte, string, string) (*nats.Msg, error), msg *nats.Msg) *natsMessage {
	return &natsMessage{Msg: msg, codec: codec, make: maker, msgCounter: msgCounter, bytesCounter: bytesCounter}
}

func (m natsMessage) Equal(msg Message) bool {
	return m.Msg.Equal(msg.Message())
}

func (m natsMessage) Message() *nats.Msg {
	return m.Msg
}

func (m natsMessage) Respond(msg proto.Message) error {
	payload, err := m.codec.Encode(nil, msg)
	if err != nil {
		return err
	}

	nmsg, err := m.make(payload, "", m.Msg.Reply)
	if err != nil {
		return err
	}

	err = m.RespondMsg(nmsg)
	m.msgCounter.Add(1)
	m.bytesCounter.Add(uint64(len(nmsg.Data)))
	return err
}

func (m natsMessage) Subject() string {
	return m.Msg.Subject
}

func (m natsMessage) Reply() string {
	return m.Msg.Reply
}

func (m natsMessage) Data() []byte {
	return m.Msg.Data
}

func (m natsMessage) Header() nats.Header {
	return m.Msg.Header
}

func (m natsMessage) QueueName() string {
	return m.Msg.Sub.Queue
}
