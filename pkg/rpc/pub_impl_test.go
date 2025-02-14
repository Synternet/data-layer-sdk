package rpc_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	nats "github.com/nats-io/nats.go"
	"github.com/synternet/data-layer-sdk/pkg/rpc"
	"github.com/synternet/data-layer-sdk/pkg/service"
	rpctypes "github.com/synternet/data-layer-sdk/types/rpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var _ rpc.Publisher = (*Publisher)(nil)

type Publisher struct {
	ctx     context.Context
	streams map[string]chan service.Message
	prefix  string
}

func NewPublisher(ctx context.Context, prefix string) *Publisher {
	return &Publisher{
		ctx:     ctx,
		streams: make(map[string]chan service.Message),
		prefix:  prefix,
	}
}

func subject(subj ...string) string {
	return strings.Join(subj, ".")
}

func replySubject(subj ...string) string {
	return "REPLY." + strings.Join(subj, ".")
}

func (p *Publisher) stream(subj ...string) chan service.Message {
	subject := subject(subj...)
	if ch, ok := p.streams[subject]; ok {
		return ch
	}
	ch := make(chan service.Message, 10)
	p.streams[subject] = ch
	return ch
}

// RequestFrom implements rpc.Publisher.
func (p *Publisher) RequestFrom(ctx context.Context, msg proto.Message, resp proto.Message, tokens ...string) (service.Message, error) {
	ch := p.stream(tokens...)
	replyCh := p.stream(replySubject(tokens...))

	msgData, err := protojson.Marshal(msg)
	if err != nil {
		panic(err)
	}

	slog.Info("requestFrom: send", "replyTo", replySubject(tokens...), "sendTo", subject(tokens...))
	select {
	case ch <- NewMsg(msgData, nil, replySubject(tokens...), subject(tokens...)):
	case <-p.ctx.Done():
		return nil, p.ctx.Err()
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-p.ctx.Done():
		return nil, p.ctx.Err()
	case data := <-replyCh:
		_, err := p.Unmarshal(data, resp)
		return data, err
	}
}

// Serve implements rpc.Publisher.
func (p *Publisher) Serve(handler service.ServiceHandler, suffixes ...string) (*nats.Subscription, error) {
	ch := p.stream(suffixes...)
	slog.Info("serve", "listenTo", subject(suffixes...))
	go func() {
		for {
			select {
			case data := <-ch:
				slog.Info("serve: received", "listenTo", subject(suffixes...), "subj", data.Subject(), "replyTo", data.Reply())
				reply, err := handler(data)
				if err != nil {
					reply = &rpctypes.Error{
						Error: err.Error(),
					}
				}
				replyCh := p.stream(data.Reply())
				msgData, err := protojson.Marshal(reply)
				if err != nil {
					panic(err)
				}

				select {
				case replyCh <- NewMsg(msgData, nil, "", data.Reply()):
				case <-p.ctx.Done():
					return
				}

			case <-p.ctx.Done():
				return
			}
		}
	}()
	return &nats.Subscription{}, nil
}

// SubscribeTo implements rpc.Publisher.
func (p *Publisher) SubscribeTo(handler service.MessageHandler, tokens ...string) (*nats.Subscription, error) {
	ch := p.stream(tokens...)
	slog.Info("subscribeTo", "listenTo", subject(tokens...))
	go func() {
		for {
			select {
			case data := <-ch:
				slog.Info("subscribeTo: received", "listenTo", subject(tokens...), "subj", data.Subject(), "replyTo", data.Reply())
				handler(data)
			case <-p.ctx.Done():
				return
			}
		}
	}()
	return &nats.Subscription{}, nil
}

func (p *Publisher) PublishTo(msg proto.Message, tokens ...string) error {
	ch := p.stream(tokens...)

	msgData, err := protojson.Marshal(msg)
	if err != nil {
		panic(err)
	}

	slog.Info("publishTo: send", "replyTo", replySubject(tokens...), "sendTo", subject(tokens...))
	select {
	case ch <- NewMsg(msgData, nil, replySubject(tokens...), subject(tokens...)):
	case <-p.ctx.Done():
		return p.ctx.Err()
	}

	return nil
}

// Unmarshal implements rpc.Publisher.
func (p *Publisher) Unmarshal(nmsg service.Message, msg proto.Message) (nats.Header, error) {
	err := json.Unmarshal(nmsg.Data(), msg)
	return nil, err
}

func (p *Publisher) Subject(suffixes ...string) string {
	tokens := []string{p.prefix}
	if suffixes != nil {
		tokens = append(tokens, suffixes...)
	}
	return strings.Join(tokens, ".")
}
