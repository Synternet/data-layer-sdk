// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: synternet/telemetry/ping.proto

package telemetry

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	PingService_PingPong_FullMethodName = "/synternet.telemetry.PingService/PingPong"
)

// PingServiceClient is the client API for PingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// PingService allows measuring latency by obtaining message flight time and recording the timestamps.
// The resulting subject is {pub_prefix}.{pub_name}.telemetry.ping
type PingServiceClient interface {
	PingPong(ctx context.Context, in *Ping, opts ...grpc.CallOption) (*Pong, error)
}

type pingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPingServiceClient(cc grpc.ClientConnInterface) PingServiceClient {
	return &pingServiceClient{cc}
}

func (c *pingServiceClient) PingPong(ctx context.Context, in *Ping, opts ...grpc.CallOption) (*Pong, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Pong)
	err := c.cc.Invoke(ctx, PingService_PingPong_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PingServiceServer is the server API for PingService service.
// All implementations should embed UnimplementedPingServiceServer
// for forward compatibility.
//
// PingService allows measuring latency by obtaining message flight time and recording the timestamps.
// The resulting subject is {pub_prefix}.{pub_name}.telemetry.ping
type PingServiceServer interface {
	PingPong(context.Context, *Ping) (*Pong, error)
}

// UnimplementedPingServiceServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPingServiceServer struct{}

func (UnimplementedPingServiceServer) PingPong(context.Context, *Ping) (*Pong, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingPong not implemented")
}
func (UnimplementedPingServiceServer) testEmbeddedByValue() {}

// UnsafePingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PingServiceServer will
// result in compilation errors.
type UnsafePingServiceServer interface {
	mustEmbedUnimplementedPingServiceServer()
}

func RegisterPingServiceServer(s grpc.ServiceRegistrar, srv PingServiceServer) {
	// If the following call pancis, it indicates UnimplementedPingServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PingService_ServiceDesc, srv)
}

func _PingService_PingPong_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PingServiceServer).PingPong(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PingService_PingPong_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PingServiceServer).PingPong(ctx, req.(*Ping))
	}
	return interceptor(ctx, in, info, handler)
}

// PingService_ServiceDesc is the grpc.ServiceDesc for PingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "synternet.telemetry.PingService",
	HandlerType: (*PingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PingPong",
			Handler:    _PingService_PingPong_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "synternet/telemetry/ping.proto",
}
