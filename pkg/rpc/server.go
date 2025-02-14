package rpc

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/synternet/data-layer-sdk/pkg/service"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	_ "github.com/synternet/data-layer-sdk/types/rpc"
)

type serviceKey string

const (
	HeaderContextKey  serviceKey = "headers"
	SubjectContextKey serviceKey = "subject"
)

func GetHeaders(ctx context.Context) (nats.Header, bool) {
	h, ok := ctx.Value(HeaderContextKey).(nats.Header)
	return h, ok
}
func GetSubject(ctx context.Context) (service.Subject, bool) {
	h, ok := ctx.Value(SubjectContextKey).(service.Subject)
	return h, ok
}

// addHeaders returns a new context with the headers added.
func addHeaders(ctx context.Context, headers nats.Header) context.Context {
	return context.WithValue(ctx, HeaderContextKey, headers)
}

// addSubject returns a new context with the subject added.
func addSubject(ctx context.Context, subject service.Subject) context.Context {
	return context.WithValue(ctx, SubjectContextKey, subject)
}

// serviceInfo wraps information about a service. It is very similar to ServiceDesc and is constructed from it for internal purposes.
type serviceInfo struct {
	// Contains the implementation for the methods in this service.
	serviceDesc *grpc.ServiceDesc
	serviceImpl any
	methods     map[string]*grpc.MethodDesc
	streams     map[string]*grpc.StreamDesc
	mdata       any
}

// ServiceRegistrar implements grpc.ServiceRegistrar for NATS
type ServiceRegistrar struct {
	mu       sync.Mutex
	group    *errgroup.Group
	pub      Publisher
	services map[string]*serviceInfo
	prefix   string
}

func NewServiceRegistrar(group *errgroup.Group, pub Publisher) *ServiceRegistrar {
	ret := &ServiceRegistrar{
		pub:      pub,
		group:    group,
		services: make(map[string]*serviceInfo),
		prefix:   "",
	}

	// pub.AddStatusCallback(ret.getStatus)
	return ret
}

func (s *ServiceRegistrar) getStatus() map[string]string {
	return map[string]string{}
}

func (s *ServiceRegistrar) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if impl != nil {
		ht := reflect.TypeOf(desc.HandlerType).Elem()
		st := reflect.TypeOf(impl)
		if !st.Implements(ht) {
			panic(fmt.Errorf("grpc: Server.RegisterService found the handler of type %v that does not satisfy %v", st, ht))
		}
	}
	s.register(desc, impl)
}

func (s *ServiceRegistrar) register(sd *grpc.ServiceDesc, ss interface{}) {
	if _, ok := s.services[sd.ServiceName]; ok {
		panic(fmt.Errorf("grpc: Server.RegisterService found duplicate service registration for %q", sd.ServiceName))
	}
	info := &serviceInfo{
		serviceImpl: ss,
		serviceDesc: sd,
		methods:     make(map[string]*grpc.MethodDesc),
		streams:     make(map[string]*grpc.StreamDesc),
		mdata:       sd.Metadata,
	}
	for i := range sd.Methods {
		d := &sd.Methods[i]
		info.methods[d.MethodName] = d
	}
	for i := range sd.Streams {
		d := &sd.Streams[i]
		info.streams[d.StreamName] = d
	}
	s.services[sd.ServiceName] = info
}

func (s *ServiceRegistrar) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, svc := range s.services {
		err := s.startUnaryMethods(ctx, svc)
		if err != nil {
			return fmt.Errorf("unary methods: %w", err)
		}
		err = s.startStreamMethods(ctx, svc)
		if err != nil {
			return fmt.Errorf("stream methods: %w", err)
		}
	}
	return nil
}

func (s *ServiceRegistrar) startUnaryMethods(ctx context.Context, svc *serviceInfo) error {
	serviceDesc, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(svc.serviceDesc.ServiceName))
	if err != nil {
		return fmt.Errorf("service desc: %w", err)
	}
	svcDesc := serviceDesc.(protoreflect.ServiceDescriptor)

	for _, method := range svc.methods {
		methodDesc := svcDesc.Methods().ByName(protoreflect.Name(method.MethodName))
		if methodDesc == nil {
			return fmt.Errorf("service desc nil: %s@%s", method.MethodName, svc.serviceDesc.ServiceName)
		}

		// deriveSubject uses svc.desc and mDesc to compute the subject tokens.
		tokens := deriveSubject(s.prefix, serviceDesc, methodDesc)
		if tokens == nil {
			return fmt.Errorf("failed to derive subject for method %s", method.MethodName)
		}
		slog.Debug("ServiceRegistrar", "service", svcDesc.FullName(), "method", method.MethodName, "subject", s.pub.Subject(tokens...))
		if skipSubscription(methodDesc) {
			continue
		}
		// Create a handler that reuses a decoder function.
		handler := func(msg service.Message) (proto.Message, error) {
			// Build a decoder function: the generated handler will call dec with a new request instance.
			dec := func(v interface{}) error {
				// We expect v to be a proto.Message.
				pm, ok := v.(proto.Message)
				if !ok {
					return fmt.Errorf("expected proto.Message, got %T", v)
				}
				_, err := s.pub.Unmarshal(msg, pm)
				return err
			}

			ctx := addSubject(ctx, service.Subject(msg.Subject()))
			ctx = addHeaders(ctx, msg.Header())

			// Invoke the generated handler directly.
			// The handler has signature:
			//   func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error)
			out, err := method.Handler(svc.serviceImpl, ctx, dec, nil)
			if err != nil {
				return nil, fmt.Errorf("%s.%s: %w", svcDesc.FullName(), method.MethodName, err)
			}
			return out.(proto.Message), nil
		}

		if _, err := s.pub.Serve(handler, tokens...); err != nil {
			return fmt.Errorf("failed to serve on %v: %w", tokens, err)
		}
	}
	return nil
}

func (s *ServiceRegistrar) startStreamMethods(ctx context.Context, svc *serviceInfo) error {
	serviceDesc, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(svc.serviceDesc.ServiceName))
	if err != nil {
		return fmt.Errorf("service desc: %w", err)
	}
	svcDesc := serviceDesc.(protoreflect.ServiceDescriptor)

	for _, stream := range svc.streams {
		methodDesc := svcDesc.Methods().ByName(protoreflect.Name(stream.StreamName))
		tokens := deriveSubject(s.prefix, svcDesc, methodDesc)
		slog.Debug("ServiceRegistrar", "service", svcDesc.FullName(), "stream", stream.StreamName, "subject", s.pub.Subject(tokens...))

		if skipSubscription(methodDesc) {
			continue
		}

		handler := func(msg service.Message) {
			ctx := addSubject(ctx, service.Subject(msg.Subject()))
			ctx = addHeaders(ctx, msg.Header())
			// Create the stream
			serverStream := &serverStream{msg: msg, ctx: ctx, pub: s.pub}

			// Call the implementation
			// type StreamHandler func(srv any, stream ServerStream) error
			if err := stream.Handler(svc.serviceImpl, serverStream); err != nil {
				msg.Nak()
				return
			}

			msg.Ack()
		}

		if _, err := s.pub.SubscribeTo(handler, s.pub.Subject(tokens...)); err != nil {
			return fmt.Errorf("failed to subscribe to %v: %v", tokens, err)
		}
	}
	return nil
}

var _ grpc.ServerStream = (*serverStream)(nil)

// serverStream implements grpc.ServerStream
type serverStream struct {
	ctx context.Context
	pub Publisher
	msg service.Message
}

func (s *serverStream) SetHeader(md metadata.MD) error {
	header := nats.Header{}
	for k, v := range md {
		header[k] = v
	}
	s.msg.Message().Header = header
	return nil
}

func (s *serverStream) SendHeader(md metadata.MD) error {
	return s.SetHeader(md)
}

func (s *serverStream) SetTrailer(md metadata.MD) {
	// No direct trailer support in NATS
}

func (s *serverStream) Context() context.Context {
	return s.ctx
}

func (s *serverStream) SendMsg(m interface{}) error {
	msg, ok := m.(proto.Message)
	if !ok {
		return fmt.Errorf("invalid message type")
	}
	return s.msg.Respond(msg)
}

func (s *serverStream) RecvMsg(m interface{}) error {
	msg, ok := m.(proto.Message)
	if !ok {
		return fmt.Errorf("invalid message type")
	}
	_, err := s.pub.Unmarshal(s.msg, msg)
	return err
}
