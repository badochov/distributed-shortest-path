// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: proto/link.proto

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

// LinkClient is the client API for Link service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LinkClient interface {
	Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddResponse, error)
	Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*InitResponse, error)
	Step(ctx context.Context, in *StepRequest, opts ...grpc.CallOption) (*StepResponse, error)
	Reconstruct(ctx context.Context, in *ReconstructRequest, opts ...grpc.CallOption) (*ReconstructResponse, error)
	Finish(ctx context.Context, in *FinishRequest, opts ...grpc.CallOption) (*FinishResponse, error)
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

func (c *linkClient) Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*InitResponse, error) {
	out := new(InitResponse)
	err := c.cc.Invoke(ctx, "/Link/Init", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *linkClient) Step(ctx context.Context, in *StepRequest, opts ...grpc.CallOption) (*StepResponse, error) {
	out := new(StepResponse)
	err := c.cc.Invoke(ctx, "/Link/Step", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *linkClient) Reconstruct(ctx context.Context, in *ReconstructRequest, opts ...grpc.CallOption) (*ReconstructResponse, error) {
	out := new(ReconstructResponse)
	err := c.cc.Invoke(ctx, "/Link/Reconstruct", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *linkClient) Finish(ctx context.Context, in *FinishRequest, opts ...grpc.CallOption) (*FinishResponse, error) {
	out := new(FinishResponse)
	err := c.cc.Invoke(ctx, "/Link/Finish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LinkServer is the server API for Link service.
// All implementations must embed UnimplementedLinkServer
// for forward compatibility
type LinkServer interface {
	Add(context.Context, *AddRequest) (*AddResponse, error)
	Init(context.Context, *InitRequest) (*InitResponse, error)
	Step(context.Context, *StepRequest) (*StepResponse, error)
	Reconstruct(context.Context, *ReconstructRequest) (*ReconstructResponse, error)
	Finish(context.Context, *FinishRequest) (*FinishResponse, error)
	mustEmbedUnimplementedLinkServer()
}

// UnimplementedLinkServer must be embedded to have forward compatible implementations.
type UnimplementedLinkServer struct {
}

func (UnimplementedLinkServer) Add(context.Context, *AddRequest) (*AddResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (UnimplementedLinkServer) Init(context.Context, *InitRequest) (*InitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Init not implemented")
}
func (UnimplementedLinkServer) Step(context.Context, *StepRequest) (*StepResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Step not implemented")
}
func (UnimplementedLinkServer) Reconstruct(context.Context, *ReconstructRequest) (*ReconstructResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reconstruct not implemented")
}
func (UnimplementedLinkServer) Finish(context.Context, *FinishRequest) (*FinishResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Finish not implemented")
}
func (UnimplementedLinkServer) mustEmbedUnimplementedLinkServer() {}

// UnsafeLinkServer may be embedded to opt out of forward compatibility for this service.
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

func _Link_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LinkServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Link/Init",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LinkServer).Init(ctx, req.(*InitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Link_Step_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StepRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LinkServer).Step(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Link/Step",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LinkServer).Step(ctx, req.(*StepRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Link_Reconstruct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReconstructRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LinkServer).Reconstruct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Link/Reconstruct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LinkServer).Reconstruct(ctx, req.(*ReconstructRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Link_Finish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FinishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LinkServer).Finish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Link/Finish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LinkServer).Finish(ctx, req.(*FinishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Link_ServiceDesc is the grpc.ServiceDesc for Link service.
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
		{
			MethodName: "Init",
			Handler:    _Link_Init_Handler,
		},
		{
			MethodName: "Step",
			Handler:    _Link_Step_Handler,
		},
		{
			MethodName: "Reconstruct",
			Handler:    _Link_Reconstruct_Handler,
		},
		{
			MethodName: "Finish",
			Handler:    _Link_Finish_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/link.proto",
}
