package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

var ErrNotAvailable = fmt.Errorf("Not available")

func (b *Service) jsStreamName(N int) string {
	if N <= 0 {
		N = len(b.streams) + 1
	}
	return fmt.Sprintf("%s-stream-%d", b.Name, N)
}

func (b *Service) jsConsumerName(N int) string {
	return fmt.Sprintf("%s-%s", b.Identity, b.Name)
}

// AddStream is an experimental feature that creates JetStream stream and durable consumer.
// The interface for this feature is experimental and is should be expected to change.
//
// Stream names must be explicit(no pattern matching) and must belong to only one stream.
//
// NOTE: Message Handler should manually acknowledge messages.
func (b *Service) AddStream(maxMsgs, maxBytes uint64, age time.Duration, subjects ...string) error {
	if b.js == nil {
		return ErrNotAvailable
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	for _, s := range subjects {
		if strings.ContainsAny(s, ">*") {
			return fmt.Errorf("%s unsupported", s)
		}

		if _, ok := b.streamSubjects[s]; ok {
			return fmt.Errorf("%s already configured", s)
		}
	}

	N := len(b.streams)
	name := b.jsStreamName(N)

	connCtx, connCancelFn := context.WithTimeout(b.Context, 10*time.Second)
	defer connCancelFn()

	cfg := &nats.StreamConfig{
		Name:     name,
		Subjects: subjects,
		MaxMsgs:  int64(maxMsgs),
		MaxBytes: int64(maxBytes),
		MaxAge:   age,
	}
	si, err := b.js.AddStream(cfg, nats.Context(connCtx))
	if err != nil {
		return fmt.Errorf("AddStream failed: %w", err)
	}

	ci, err := b.js.AddConsumer(name, &nats.ConsumerConfig{
		Durable: b.jsConsumerName(N),
	})
	if err != nil {
		return fmt.Errorf("AddConsumer failed: %w", err)
	}

	b.streams = append(b.streams, jsStream{
		cfg:          cfg,
		streamInfo:   si,
		consumerInfo: ci,
	})

	for _, s := range subjects {
		b.streamSubjects[s] = N
	}

	return nil
}

func (b *Service) attemptJSConsume(handler nats.MsgHandler, subject string) (*nats.Subscription, error) {
	if b.js == nil {
		return nil, ErrNotAvailable
	}

	streamName, consumerName := "", ""
	if id, ok := b.streamSubjects[subject]; !ok {
		return nil, ErrNotAvailable
	} else {
		info := b.streams[id]
		streamName = info.cfg.Name
		consumerName = info.consumerInfo.Config.Durable
	}

	return b.js.Subscribe(subject, handler, nats.Bind(streamName, consumerName), nats.ManualAck())
}
