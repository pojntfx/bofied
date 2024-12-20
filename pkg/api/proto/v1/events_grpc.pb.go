// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.19.6
// source: events.proto

package v1

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
	EventsService_SubscribeToEvents_FullMethodName = "/com.pojtinger.felicitas.bofied.EventsService/SubscribeToEvents"
)

// EventsServiceClient is the client API for EventsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventsServiceClient interface {
	SubscribeToEvents(ctx context.Context, in *Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[EventMessage], error)
}

type eventsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEventsServiceClient(cc grpc.ClientConnInterface) EventsServiceClient {
	return &eventsServiceClient{cc}
}

func (c *eventsServiceClient) SubscribeToEvents(ctx context.Context, in *Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[EventMessage], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &EventsService_ServiceDesc.Streams[0], EventsService_SubscribeToEvents_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[Empty, EventMessage]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type EventsService_SubscribeToEventsClient = grpc.ServerStreamingClient[EventMessage]

// EventsServiceServer is the server API for EventsService service.
// All implementations must embed UnimplementedEventsServiceServer
// for forward compatibility.
type EventsServiceServer interface {
	SubscribeToEvents(*Empty, grpc.ServerStreamingServer[EventMessage]) error
	mustEmbedUnimplementedEventsServiceServer()
}

// UnimplementedEventsServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEventsServiceServer struct{}

func (UnimplementedEventsServiceServer) SubscribeToEvents(*Empty, grpc.ServerStreamingServer[EventMessage]) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeToEvents not implemented")
}
func (UnimplementedEventsServiceServer) mustEmbedUnimplementedEventsServiceServer() {}
func (UnimplementedEventsServiceServer) testEmbeddedByValue()                       {}

// UnsafeEventsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventsServiceServer will
// result in compilation errors.
type UnsafeEventsServiceServer interface {
	mustEmbedUnimplementedEventsServiceServer()
}

func RegisterEventsServiceServer(s grpc.ServiceRegistrar, srv EventsServiceServer) {
	// If the following call pancis, it indicates UnimplementedEventsServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&EventsService_ServiceDesc, srv)
}

func _EventsService_SubscribeToEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EventsServiceServer).SubscribeToEvents(m, &grpc.GenericServerStream[Empty, EventMessage]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type EventsService_SubscribeToEventsServer = grpc.ServerStreamingServer[EventMessage]

// EventsService_ServiceDesc is the grpc.ServiceDesc for EventsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.pojtinger.felicitas.bofied.EventsService",
	HandlerType: (*EventsServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeToEvents",
			Handler:       _EventsService_SubscribeToEvents_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "events.proto",
}
