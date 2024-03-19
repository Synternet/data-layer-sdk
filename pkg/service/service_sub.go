package service

import (
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/syntropynet/data-layer-sdk/pkg/options"
)

// Serve is a convenience method to serve a service subject. It acts the same as Subscribe, but takes `ServiceHandler` instead and will respond
// either with Error type, or response from the handler. Serve will use ReqNats connection.
func (b *Service) Serve(handler ServiceHandler, suffixes ...string) (*nats.Subscription, error) {
	return b.subscribeTo(
		b.ReqNats,
		func(msg Message) {
			resp, err := handler(msg)
			if err != nil {
				b.Logger.Error("service handler failed", "err", err, "suffixes", suffixes)
				err1 := msg.Respond(&Error{Error: err.Error()})
				if err1 != nil {
					b.Logger.Error("service handler failed during error", "err", err, "err1", err1, "suffixes", suffixes)
				}
				return
			}
			err = msg.Respond(resp)
			if err != nil {
				b.Logger.Error("service handler failed", "err", err, "suffixes", suffixes)
			}
		},
		b.Subject(suffixes...),
	)
}

// Subscribe will subscribe to a subject constructed from {prefix}.{name}.{...suffixes}, where
// suffixes are joined using ".". Subscribe will use SubNats connection.
func (b *Service) Subscribe(handler MessageHandler, suffixes ...string) (*nats.Subscription, error) {
	return b.subscribeTo(b.SubNats, handler, b.Subject(suffixes...))
}

// SubscribeTo will subscribe to a subject constructed {...suffixes}, where
// suffixes are joined using ".".
//
// Experimental: When a stream was registered with AddStream SubscribeTo will use durable stream instead of realtime.
func (b *Service) SubscribeTo(handler MessageHandler, suffixes ...string) (*nats.Subscription, error) {
	return b.subscribeTo(b.SubNats, handler, suffixes...)
}

func (b *Service) subscribeTo(nc options.NatsConn, handler MessageHandler, suffixes ...string) (*nats.Subscription, error) {
	if nc == nil {
		return nil, fmt.Errorf("subscribing NATS connection is nil")
	}
	if b.VerboseLog {
		b.Logger.Debug("SubscribeTo", "suffixes", suffixes)
	}
	natsHandler := func(msg *nats.Msg) {
		b.msg_in_counter.Add(1)
		b.bytes_in_counter.Add(uint64(len(msg.Data)))
		wrapped := wrapMessage(b.Codec, &b.msg_out_counter, &b.bytes_out_counter, b.makeMsg, msg)
		handler(wrapped)
	}

	subject := strings.Join(suffixes, ".")

	if sub, err := b.attemptJSConsume(natsHandler, subject); err == nil {
		return sub, nil
	}

	if b.QueueName != "" {
		return nc.QueueSubscribe(subject, b.QueueName, natsHandler)
	}
	return nc.Subscribe(subject, natsHandler)
}
