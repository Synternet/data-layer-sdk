package service

import (
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/syntropynet/data-layer-sdk/pkg/options"
)

// Subscribe will subscribe to a subject constructed from {prefix}.{name}.{...suffixes}, where
// suffixes are joined using ".". Subscribe will use SubNats connection.
func (b *Service) Subscribe(handler MessageHandler, suffixes ...string) (*nats.Subscription, error) {
	return b.subscribeTo(b.SubNats, handler, b.Subject(suffixes...))
}

// SubscribeTo will subscribe to a subject constructed as {...tokens}, where
// tokens are joined using ".".
//
// Experimental: When a stream was registered with AddStream SubscribeTo will use durable stream instead of realtime.
func (b *Service) SubscribeTo(handler MessageHandler, tokens ...string) (*nats.Subscription, error) {
	return b.subscribeTo(b.SubNats, handler, tokens...)
}

func (b *Service) subscribeTo(nc options.NatsConn, handler MessageHandler, tokens ...string) (*nats.Subscription, error) {
	if nc == nil {
		return nil, ErrSubConnection
	}
	if b.VerboseLog {
		b.Logger.Debug("subscribeTo", "tokens", tokens)
	}
	natsHandler := func(msg *nats.Msg) {
		b.msg_in_counter.Add(1)
		b.bytes_in_counter.Add(uint64(len(msg.Data)))
		wrapped := wrapMessage(b.Codec, &b.msg_out_counter, &b.bytes_out_counter, b.makeMsg, msg)
		handler(wrapped)
	}

	subject := strings.Join(tokens, ".")

	if sub, err := b.attemptJSConsume(natsHandler, subject); err == nil {
		return sub, nil
	}

	if b.QueueName != "" {
		return nc.QueueSubscribe(subject, b.QueueName, natsHandler)
	}
	return nc.Subscribe(subject, natsHandler)
}
