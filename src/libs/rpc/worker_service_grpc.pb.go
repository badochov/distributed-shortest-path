// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: worker_service.proto

package rpc

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

// WorkerClient is the client API for Worker service_server.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WorkerClient interface {
	AssignSegment(ctx context.Context, in *Segment, opts ...grpc.CallOption) (*Ack, error)
}

type workerClient struct {
	cc grpc.ClientConnInterface
}

func NewWorkerClient(cc grpc.ClientConnInterface) WorkerClient {
	return &workerClient{cc}
}

func (c *workerClient) AssignSegment(ctx context.Context, in *Segment, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/Worker/AssignSegment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WorkerServer is the server API for Worker service_server.
// All implementations must embed UnimplementedWorkerServer
// for forward compatibility
type WorkerServer interface {
	AssignSegment(context.Context, *Segment) (*Ack, error)
	mustEmbedUnimplementedWorkerServer()
}

// UnimplementedWorkerServer must be embedded to have forward compatible implementations.
type UnimplementedWorkerServer struct {
}

func (UnimplementedWorkerServer) AssignSegment(context.Context, *Segment) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignSegment not implemented")
}
func (UnimplementedWorkerServer) mustEmbedUnimplementedWorkerServer() {}

// UnsafeWorkerServer may be embedded to opt out of forward compatibility for this service_server.
// Use of this interface is not recommended, as added methods to WorkerServer will
// result in compilation errors.
type UnsafeWorkerServer interface {
	mustEmbedUnimplementedWorkerServer()
}

func RegisterWorkerServer(s grpc.ServiceRegistrar, srv WorkerServer) {
	s.RegisterService(&Worker_ServiceDesc, srv)
}

func _Worker_AssignSegment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Segment)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkerServer).AssignSegment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Worker/AssignSegment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkerServer).AssignSegment(ctx, req.(*Segment))
	}
	return interceptor(ctx, in, info, handler)
}

// Worker_ServiceDesc is the grpc.ServiceDesc for Worker service_server.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Worker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Worker",
	HandlerType: (*WorkerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AssignSegment",
			Handler:    _Worker_AssignSegment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "worker_service.proto",
}
