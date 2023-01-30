// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: proto/link.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AddRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	A int32 `protobuf:"varint,1,opt,name=a,proto3" json:"a,omitempty"`
	B int32 `protobuf:"varint,2,opt,name=b,proto3" json:"b,omitempty"`
}

func (x *AddRequest) Reset() {
	*x = AddRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_link_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddRequest) ProtoMessage() {}

func (x *AddRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_link_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddRequest.ProtoReflect.Descriptor instead.
func (*AddRequest) Descriptor() ([]byte, []int) {
	return file_proto_link_proto_rawDescGZIP(), []int{0}
}

func (x *AddRequest) GetA() int32 {
	if x != nil {
		return x.A
	}
	return 0
}

func (x *AddRequest) GetB() int32 {
	if x != nil {
		return x.B
	}
	return 0
}

type AddResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Res int32 `protobuf:"varint,1,opt,name=res,proto3" json:"res,omitempty"`
}

func (x *AddResponse) Reset() {
	*x = AddResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_link_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddResponse) ProtoMessage() {}

func (x *AddResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_link_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddResponse.ProtoReflect.Descriptor instead.
func (*AddResponse) Descriptor() ([]byte, []int) {
	return file_proto_link_proto_rawDescGZIP(), []int{1}
}

func (x *AddResponse) GetRes() int32 {
	if x != nil {
		return x.Res
	}
	return 0
}

type InitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestId uint64 `protobuf:"varint,1,opt,name=requestId,proto3" json:"requestId,omitempty"`
}

func (x *InitRequest) Reset() {
	*x = InitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_link_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitRequest) ProtoMessage() {}

func (x *InitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_link_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitRequest.ProtoReflect.Descriptor instead.
func (*InitRequest) Descriptor() ([]byte, []int) {
	return file_proto_link_proto_rawDescGZIP(), []int{2}
}

func (x *InitRequest) GetRequestId() uint64 {
	if x != nil {
		return x.RequestId
	}
	return 0
}

type InitResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InitResponse) Reset() {
	*x = InitResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_link_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitResponse) ProtoMessage() {}

func (x *InitResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_link_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitResponse.ProtoReflect.Descriptor instead.
func (*InitResponse) Descriptor() ([]byte, []int) {
	return file_proto_link_proto_rawDescGZIP(), []int{3}
}

type StepRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VertexId  int64   `protobuf:"varint,1,opt,name=vertexId,proto3" json:"vertexId,omitempty"`
	Distance  float64 `protobuf:"fixed64,2,opt,name=distance,proto3" json:"distance,omitempty"`
	Through   int64   `protobuf:"varint,3,opt,name=through,proto3" json:"through,omitempty"`
	RequestId uint64  `protobuf:"varint,4,opt,name=requestId,proto3" json:"requestId,omitempty"`
}

func (x *StepRequest) Reset() {
	*x = StepRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_link_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StepRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StepRequest) ProtoMessage() {}

func (x *StepRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_link_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StepRequest.ProtoReflect.Descriptor instead.
func (*StepRequest) Descriptor() ([]byte, []int) {
	return file_proto_link_proto_rawDescGZIP(), []int{4}
}

func (x *StepRequest) GetVertexId() int64 {
	if x != nil {
		return x.VertexId
	}
	return 0
}

func (x *StepRequest) GetDistance() float64 {
	if x != nil {
		return x.Distance
	}
	return 0
}

func (x *StepRequest) GetThrough() int64 {
	if x != nil {
		return x.Through
	}
	return 0
}

func (x *StepRequest) GetRequestId() uint64 {
	if x != nil {
		return x.RequestId
	}
	return 0
}

type StepResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VertexId int64   `protobuf:"varint,1,opt,name=vertexId,proto3" json:"vertexId,omitempty"`
	Distance float64 `protobuf:"fixed64,2,opt,name=distance,proto3" json:"distance,omitempty"`
	Through  int64   `protobuf:"varint,3,opt,name=through,proto3" json:"through,omitempty"`
}

func (x *StepResponse) Reset() {
	*x = StepResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_link_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StepResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StepResponse) ProtoMessage() {}

func (x *StepResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_link_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StepResponse.ProtoReflect.Descriptor instead.
func (*StepResponse) Descriptor() ([]byte, []int) {
	return file_proto_link_proto_rawDescGZIP(), []int{5}
}

func (x *StepResponse) GetVertexId() int64 {
	if x != nil {
		return x.VertexId
	}
	return 0
}

func (x *StepResponse) GetDistance() float64 {
	if x != nil {
		return x.Distance
	}
	return 0
}

func (x *StepResponse) GetThrough() int64 {
	if x != nil {
		return x.Through
	}
	return 0
}

type ReconstructRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VertexId  int64  `protobuf:"varint,1,opt,name=vertexId,proto3" json:"vertexId,omitempty"`
	RequestId uint64 `protobuf:"varint,2,opt,name=requestId,proto3" json:"requestId,omitempty"`
}

func (x *ReconstructRequest) Reset() {
	*x = ReconstructRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_link_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReconstructRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReconstructRequest) ProtoMessage() {}

func (x *ReconstructRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_link_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReconstructRequest.ProtoReflect.Descriptor instead.
func (*ReconstructRequest) Descriptor() ([]byte, []int) {
	return file_proto_link_proto_rawDescGZIP(), []int{6}
}

func (x *ReconstructRequest) GetVertexId() int64 {
	if x != nil {
		return x.VertexId
	}
	return 0
}

func (x *ReconstructRequest) GetRequestId() uint64 {
	if x != nil {
		return x.RequestId
	}
	return 0
}

type ReconstructResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path []int64 `protobuf:"varint,1,rep,packed,name=path,proto3" json:"path,omitempty"`
}

func (x *ReconstructResponse) Reset() {
	*x = ReconstructResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_link_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReconstructResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReconstructResponse) ProtoMessage() {}

func (x *ReconstructResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_link_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReconstructResponse.ProtoReflect.Descriptor instead.
func (*ReconstructResponse) Descriptor() ([]byte, []int) {
	return file_proto_link_proto_rawDescGZIP(), []int{7}
}

func (x *ReconstructResponse) GetPath() []int64 {
	if x != nil {
		return x.Path
	}
	return nil
}

var File_proto_link_proto protoreflect.FileDescriptor

var file_proto_link_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6c, 0x69, 0x6e, 0x6b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x28, 0x0a, 0x0a, 0x41, 0x64, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0c, 0x0a, 0x01, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x61, 0x12, 0x0c,
	0x0a, 0x01, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x62, 0x22, 0x1f, 0x0a, 0x0b,
	0x41, 0x64, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x72,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x72, 0x65, 0x73, 0x22, 0x2b, 0x0a,
	0x0b, 0x49, 0x6e, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x22, 0x0e, 0x0a, 0x0c, 0x49, 0x6e,
	0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x7d, 0x0a, 0x0b, 0x53, 0x74,
	0x65, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x76, 0x65, 0x72,
	0x74, 0x65, 0x78, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x76, 0x65, 0x72,
	0x74, 0x65, 0x78, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x68, 0x72, 0x6f, 0x75, 0x67, 0x68, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x07, 0x74, 0x68, 0x72, 0x6f, 0x75, 0x67, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x22, 0x60, 0x0a, 0x0c, 0x53, 0x74, 0x65,
	0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x76, 0x65, 0x72,
	0x74, 0x65, 0x78, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x76, 0x65, 0x72,
	0x74, 0x65, 0x78, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x68, 0x72, 0x6f, 0x75, 0x67, 0x68, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x07, 0x74, 0x68, 0x72, 0x6f, 0x75, 0x67, 0x68, 0x22, 0x4e, 0x0a, 0x12, 0x52,
	0x65, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1a, 0x0a, 0x08, 0x76, 0x65, 0x72, 0x74, 0x65, 0x78, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x08, 0x76, 0x65, 0x72, 0x74, 0x65, 0x78, 0x49, 0x64, 0x12, 0x1c, 0x0a,
	0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x22, 0x29, 0x0a, 0x13, 0x52,
	0x65, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03,
	0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x32, 0xb4, 0x01, 0x0a, 0x04, 0x4c, 0x69, 0x6e, 0x6b, 0x12,
	0x22, 0x0a, 0x03, 0x41, 0x64, 0x64, 0x12, 0x0b, 0x2e, 0x41, 0x64, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x41, 0x64, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x25, 0x0a, 0x04, 0x49, 0x6e, 0x69, 0x74, 0x12, 0x0c, 0x2e, 0x49, 0x6e,
	0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x49, 0x6e, 0x69, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x25, 0x0a, 0x04, 0x53, 0x74,
	0x65, 0x70, 0x12, 0x0c, 0x2e, 0x53, 0x74, 0x65, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x0d, 0x2e, 0x53, 0x74, 0x65, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x3a, 0x0a, 0x0b, 0x52, 0x65, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74,
	0x12, 0x13, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72,
	0x75, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x4e, 0x5a,
	0x4c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x61, 0x64, 0x6f,
	0x63, 0x68, 0x6f, 0x76, 0x2f, 0x64, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64,
	0x2d, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x73, 0x74, 0x2d, 0x70, 0x61, 0x74, 0x68, 0x2f, 0x73,
	0x72, 0x63, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x77, 0x6f, 0x72, 0x6b,
	0x65, 0x72, 0x2f, 0x6c, 0x69, 0x6e, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_link_proto_rawDescOnce sync.Once
	file_proto_link_proto_rawDescData = file_proto_link_proto_rawDesc
)

func file_proto_link_proto_rawDescGZIP() []byte {
	file_proto_link_proto_rawDescOnce.Do(func() {
		file_proto_link_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_link_proto_rawDescData)
	})
	return file_proto_link_proto_rawDescData
}

var file_proto_link_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proto_link_proto_goTypes = []interface{}{
	(*AddRequest)(nil),          // 0: AddRequest
	(*AddResponse)(nil),         // 1: AddResponse
	(*InitRequest)(nil),         // 2: InitRequest
	(*InitResponse)(nil),        // 3: InitResponse
	(*StepRequest)(nil),         // 4: StepRequest
	(*StepResponse)(nil),        // 5: StepResponse
	(*ReconstructRequest)(nil),  // 6: ReconstructRequest
	(*ReconstructResponse)(nil), // 7: ReconstructResponse
}
var file_proto_link_proto_depIdxs = []int32{
	0, // 0: Link.Add:input_type -> AddRequest
	2, // 1: Link.Init:input_type -> InitRequest
	4, // 2: Link.Step:input_type -> StepRequest
	6, // 3: Link.Reconstruct:input_type -> ReconstructRequest
	1, // 4: Link.Add:output_type -> AddResponse
	3, // 5: Link.Init:output_type -> InitResponse
	5, // 6: Link.Step:output_type -> StepResponse
	7, // 7: Link.Reconstruct:output_type -> ReconstructResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_link_proto_init() }
func file_proto_link_proto_init() {
	if File_proto_link_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_link_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_link_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_link_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_link_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_link_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StepRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_link_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StepResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_link_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReconstructRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_link_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReconstructResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_link_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_link_proto_goTypes,
		DependencyIndexes: file_proto_link_proto_depIdxs,
		MessageInfos:      file_proto_link_proto_msgTypes,
	}.Build()
	File_proto_link_proto = out.File
	file_proto_link_proto_rawDesc = nil
	file_proto_link_proto_goTypes = nil
	file_proto_link_proto_depIdxs = nil
}
