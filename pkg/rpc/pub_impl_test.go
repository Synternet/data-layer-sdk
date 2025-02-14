package rpc_test

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	nats "github.com/nats-io/nats.go"
	"github.com/synternet/data-layer-sdk/pkg/rpc"
	"github.com/synternet/data-layer-sdk/pkg/service"
	rpctypes "github.com/synternet/data-layer-sdk/x/synternet/rpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var _ rpc.Publisher = (*Publisher)(nil)

type Publisher struct {
	ctx     context.Context
	t       *testing.T
	mu      sync.Mutex
	streams map[string]chan service.Message
	prefix  string
}

func NewPublisher(ctx context.Context, t *testing.T, prefix string) *Publisher {
	return &Publisher{
		ctx:     ctx,
		t:       t,
		streams: make(map[string]chan service.Message),
		prefix:  prefix,
	}
}

func subject(subj ...string) string {
	return strings.Join(subj, ".")
}

func (p *Publisher) stream(subj ...string) chan service.Message {
	p.mu.Lock()
	defer p.mu.Unlock()
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
	replyTo := p.RpcInbox(tokens...)
	replyCh := p.stream(replyTo)

	msgData, err := protojson.Marshal(msg)
	if err != nil {
		panic(err)
	}

	p.t.Log("requestFrom: send", "replyTo=", replyTo, "sendTo=", subject(tokens...), "ch=", ch, "replyCh=", replyCh)
	select {
	case <-p.ctx.Done():
		return nil, p.ctx.Err()
	case ch <- NewMsg(p.t, msgData, replyCh, replyTo, subject(tokens...)):
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-p.ctx.Done():
		return nil, p.ctx.Err()
	case data := <-replyCh:
		p.t.Log("requestFrom: received", "replyTo=", replyTo, "sendTo=", subject(tokens...))
		_, err := p.Unmarshal(data, resp)
		return data, err
	}
}

func (p *Publisher) RpcInbox(suffixes ...string) string {
	return nats.NewInbox()
}

// Serve implements rpc.Publisher.
func (p *Publisher) Serve(handler service.ServiceHandler, suffixes ...string) (*nats.Subscription, error) {
	subject := p.Subject(suffixes...)
	ch := p.stream(subject)
	p.t.Log("serve", "listenTo=", subject, "ch=", ch)
	go func() {
		for {
			select {
			case <-p.ctx.Done():
				return
			case data := <-ch:
				p.t.Log("serve: received", "listenTo=", subject, "subj=", data.Subject(), "replyTo=", data.Reply())
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
				case <-p.ctx.Done():
					return
				case replyCh <- NewMsg(p.t, msgData, nil, "", data.Reply()):
					p.t.Log("serve: publish", "subj=", data.Reply())
				}
			}
		}
	}()
	return &nats.Subscription{}, nil
}

// SubscribeTo implements rpc.Publisher.
func (p *Publisher) SubscribeTo(handler service.MessageHandler, tokens ...string) (*nats.Subscription, error) {
	ch := p.stream(tokens...)
	p.t.Log("subscribeTo", "listenTo=", subject(tokens...), "ch=", ch)
	go func() {
		for {
			select {
			case <-p.ctx.Done():
				return
			case data := <-ch:
				p.t.Log("subscribeTo: received", "listenTo=", subject(tokens...), "subj=", data.Subject(), "replyTo=", data.Reply())
				handler(data)
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

	p.t.Log("publishTo: send", "sendTo=", subject(tokens...), "ch=", ch)
	select {
	case <-p.ctx.Done():
		return p.ctx.Err()
	case ch <- NewMsg(p.t, msgData, nil, "", subject(tokens...)):
	}

	return nil
}

func (p *Publisher) PublishToRpc(msg proto.Message, replyTo string, tokens ...string) error {
	ch := p.stream(tokens...)
	replyCh := p.stream(replyTo)

	msgData, err := protojson.Marshal(msg)
	if err != nil {
		panic(err)
	}

	p.t.Log("publishToRpc: send", "replyTo=", replyTo, "sendTo=", subject(tokens...), "ch=", ch, "replyCh=", replyCh)
	select {
	case <-p.ctx.Done():
		return p.ctx.Err()
	case ch <- NewMsg(p.t, msgData, replyCh, replyTo, subject(tokens...)):
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
