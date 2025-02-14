package rpc

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/nats-io/nats.go"
	service "github.com/synternet/data-layer-sdk/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ClientConn implements grpc.ClientConnInterface for NATS
type ClientConn struct {
	ctx    context.Context
	sub    Publisher
	prefix string
}

func NewClientConn(ctx context.Context, sub Publisher, remotePrefix string) *ClientConn {
	return &ClientConn{
		sub:    sub,
		ctx:    ctx,
		prefix: remotePrefix,
	}
}

func parseServiceMethod(m string) (protoreflect.ServiceDescriptor, protoreflect.MethodDescriptor, error) {
	tmp := strings.TrimPrefix(m, "/")
	parts := strings.Split(tmp, "/")
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid method: %s", m)
	}
	service := protoreflect.FullName(parts[0])
	method := protoreflect.Name(parts[1])

	svcDescriptor, err := protoregistry.GlobalFiles.FindDescriptorByName(service)
	if err != nil {
		return nil, nil, fmt.Errorf("service descriptor: %w", err)
	}

	svcDesc := svcDescriptor.(protoreflect.ServiceDescriptor)
	return svcDesc, svcDesc.Methods().ByName(method), nil
}

func (c *ClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	svcDesc, methodDesc, err := parseServiceMethod(method)
	if err != nil {
		return fmt.Errorf("parse method: %w", err)
	}
	tokens := deriveSubject(c.prefix, svcDesc, methodDesc)
	if tokens == nil {
		return fmt.Errorf("invalid subject: %s@%s", methodDesc.FullName(), svcDesc.FullName())
	}
	slog.Debug("ClientConn.Invoke", "service", svcDesc.FullName(), "method", methodDesc.FullName(), "subject", strings.Join(tokens, "."))
	if skipSubscription(methodDesc) {
		return fmt.Errorf("calling disabled: %s@%s", methodDesc.FullName(), svcDesc.FullName())
	}
	_, err = c.sub.RequestFrom(ctx, args.(proto.Message), reply.(proto.Message), tokens...)
	return err
}

func (c *ClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	svcDesc, methodDesc, err := parseServiceMethod(method)
	if err != nil {
		return nil, fmt.Errorf("parse method: %v", err)
	}
	tokens := deriveSubject(c.prefix, svcDesc, methodDesc)
	if tokens == nil {
		return nil, fmt.Errorf("invalid subject: %s@%s", methodDesc.FullName(), svcDesc.FullName())
	}
	slog.Debug("ClientConn.NewStream", "service", svcDesc.FullName(), "method", methodDesc.FullName(), "subject", strings.Join(tokens, "."))
	if skipSubscription(methodDesc) {
		return nil, fmt.Errorf("calling disabled: %s@%s", methodDesc.FullName(), svcDesc.FullName())
	}

	stream, err := newClientStream(ctx, c.sub, tokens)
	if err != nil {
		return nil, fmt.Errorf("couldn't create client stream for %s: %w", strings.Join(tokens, "."), err)
	}

	return stream, nil
}

// clientStream implements grpc.ClientStream
type clientStream struct {
	ctx          context.Context
	pub          Publisher
	tokens       []string
	replySubject string

	mu        sync.Mutex
	sub       *nats.Subscription
	recvChan  chan service.Message
	errChan   chan error
	closeChan chan struct{}
	closed    bool
}

func newClientStream(ctx context.Context, pub Publisher, tokens []string) (*clientStream, error) {
	stream := &clientStream{
		ctx:          ctx,
		pub:          pub,
		tokens:       tokens,
		replySubject: nats.NewInbox(),
		recvChan:     make(chan service.Message, 1000),
		errChan:      make(chan error, 1),
		closeChan:    make(chan struct{}),
	}

	return stream, nil
}

func (s *clientStream) Header() (metadata.MD, error) {
	return nil, nil
}

func (s *clientStream) Trailer() metadata.MD {
	return nil
}

func (s *clientStream) CloseSend() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.closed {
		close(s.closeChan)
		s.closed = true
	}
	return nil
}

func (s *clientStream) Context() context.Context {
	return s.ctx
}

func (s *clientStream) subscribe() error {
	handler := func(msg service.Message) {
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.closed {
			return
		}

		select {
		case s.recvChan <- msg:
		case <-s.ctx.Done():
			return
		}
	}
	sub, err := s.pub.SubscribeTo(handler, s.replySubject)
	s.sub = sub
	return err
}

func (s *clientStream) SendMsg(m interface{}) error {
	if _, ok := m.(*emptypb.Empty); ok {
		return s.subscribe()
	}

	err := s.subscribe()
	if err != nil {
		return err
	}

	return s.pub.PublishTo(m.(proto.Message), s.tokens...)
}

func (s *clientStream) RecvMsg(m interface{}) error {
	if _, ok := m.(*emptypb.Empty); ok {
		return nil
	}

	select {
	case err := <-s.errChan:
		return err
	case msg := <-s.recvChan:
		_, err := s.pub.Unmarshal(msg, m.(proto.Message))
		return err
	case <-s.ctx.Done():
		return s.ctx.Err()
	case <-s.closeChan:
		return fmt.Errorf("stream closed")
	}
}

func (s *clientStream) Close() error {
	s.mu.Lock()

	defer s.mu.Unlock()

	if !s.closed {
		close(s.closeChan)
		s.closed = true
		if s.sub != nil {
			s.sub.Unsubscribe()
		}
	}
	return nil
}
