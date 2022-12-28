// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: link.proto

package proto

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

// LinkClient is the client API for Link service_server.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LinkClient interface {
	Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddResponse, error)
}

type linkClient struct {
	cc grpc.ClientConnInterface
}

func NewLinkClient(cc grpc.ClientConnInterface) LinkClient {
	return &linkClient{cc}
}

func (c *linkClient) Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddResponse, error) {
	out := new(AddResponse)
	err := c.cc.Invoke(ctx, "/Link/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LinkServer is the server API for Link service_server.
// All implementations must embed UnimplementedLinkServer
// for forward compatibility
type LinkServer interface {
	Add(context.Context, *AddRequest) (*AddResponse, error)
	mustEmbedUnimplementedLinkServer()
}

// UnimplementedLinkServer must be embedded to have forward compatible implementations.
type UnimplementedLinkServer struct {
}

func (UnimplementedLinkServer) Add(context.Context, *AddRequest) (*AddResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (UnimplementedLinkServer) mustEmbedUnimplementedLinkServer() {}

// UnsafeLinkServer may be embedded to opt out of forward compatibility for this service_server.
// Use of this interface is not recommended, as added methods to LinkServer will
// result in compilation errors.
type UnsafeLinkServer interface {
	mustEmbedUnimplementedLinkServer()
}

func RegisterLinkServer(s grpc.ServiceRegistrar, srv LinkServer) {
	s.RegisterService(&Link_ServiceDesc, srv)
}

func _Link_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LinkServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Link/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LinkServer).Add(ctx, req.(*AddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Link_ServiceDesc is the grpc.ServiceDesc for Link service_server.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Link_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Link",
	HandlerType: (*LinkServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Add",
			Handler:    _Link_Add_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "link.proto",
}
