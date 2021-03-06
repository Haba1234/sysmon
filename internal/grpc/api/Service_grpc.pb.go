// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StatisticsClient is the client API for Statistics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StatisticsClient interface {
	ListStatistics(ctx context.Context, in *SubscriptionRequest, opts ...grpc.CallOption) (Statistics_ListStatisticsClient, error)
}

type statisticsClient struct {
	cc grpc.ClientConnInterface
}

func NewStatisticsClient(cc grpc.ClientConnInterface) StatisticsClient {
	return &statisticsClient{cc}
}

func (c *statisticsClient) ListStatistics(ctx context.Context, in *SubscriptionRequest, opts ...grpc.CallOption) (Statistics_ListStatisticsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Statistics_ServiceDesc.Streams[0], "/service.Statistics/ListStatistics", opts...)
	if err != nil {
		return nil, err
	}
	x := &statisticsListStatisticsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Statistics_ListStatisticsClient interface {
	Recv() (*StatisticsResponse, error)
	grpc.ClientStream
}

type statisticsListStatisticsClient struct {
	grpc.ClientStream
}

func (x *statisticsListStatisticsClient) Recv() (*StatisticsResponse, error) {
	m := new(StatisticsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StatisticsServer is the server API for Statistics service.
// All implementations must embed UnimplementedStatisticsServer
// for forward compatibility
type StatisticsServer interface {
	ListStatistics(*SubscriptionRequest, Statistics_ListStatisticsServer) error
	mustEmbedUnimplementedStatisticsServer()
}

// UnimplementedStatisticsServer must be embedded to have forward compatible implementations.
type UnimplementedStatisticsServer struct {
}

func (UnimplementedStatisticsServer) ListStatistics(*SubscriptionRequest, Statistics_ListStatisticsServer) error {
	return status.Errorf(codes.Unimplemented, "method ListStatistics not implemented")
}
func (UnimplementedStatisticsServer) mustEmbedUnimplementedStatisticsServer() {}

// UnsafeStatisticsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StatisticsServer will
// result in compilation errors.
type UnsafeStatisticsServer interface {
	mustEmbedUnimplementedStatisticsServer()
}

func RegisterStatisticsServer(s grpc.ServiceRegistrar, srv StatisticsServer) {
	s.RegisterService(&Statistics_ServiceDesc, srv)
}

func _Statistics_ListStatistics_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscriptionRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StatisticsServer).ListStatistics(m, &statisticsListStatisticsServer{stream})
}

type Statistics_ListStatisticsServer interface {
	Send(*StatisticsResponse) error
	grpc.ServerStream
}

type statisticsListStatisticsServer struct {
	grpc.ServerStream
}

func (x *statisticsListStatisticsServer) Send(m *StatisticsResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Statistics_ServiceDesc is the grpc.ServiceDesc for Statistics service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Statistics_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.Statistics",
	HandlerType: (*StatisticsServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListStatistics",
			Handler:       _Statistics_ListStatistics_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "Service.proto",
}
