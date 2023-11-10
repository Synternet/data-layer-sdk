package options

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
)

type natsStub struct {
	Verbose *bool
}

func (s *natsStub) Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	if !*s.Verbose {
		return nil, nil
	}
	log.Println("NATS: subscribing to subj=", subj)
	return nil, nil
}

func (s *natsStub) QueueSubscribe(subj, queue string, cb nats.MsgHandler) (*nats.Subscription, error) {
	if !*s.Verbose {
		return nil, nil
	}
	log.Println("NATS: subscribing to subj=", subj, "queue=", queue)
	return nil, nil
}

func (s *natsStub) RequestMsgWithContext(ctx context.Context, m *nats.Msg) (*nats.Msg, error) {
	if !*s.Verbose {
		return m, nil
	}
	log.Println("NATS: req subj=", m.Subject, " len(data)=", len(m.Data), " data=", string(m.Data))
	return m, nil
}

func (s *natsStub) PublishMsg(m *nats.Msg) error {
	if !*s.Verbose {
		return nil
	}
	log.Println("NATS: pub subj=", m.Subject, " len(data)=", len(m.Data), " data=", string(m.Data))
	return nil
}

func (s *natsStub) Flush() error {
	if !*s.Verbose {
		return nil
	}
	return nil
}
