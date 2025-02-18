// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: example/v1/example.proto

package v1

import (
	rpc "github.com/synternet/data-layer-sdk/x/synternet/rpc"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AddRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Surname       string                 `protobuf:"bytes,2,opt,name=surname,proto3" json:"surname,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddRequest) Reset() {
	*x = AddRequest{}
	mi := &file_example_v1_example_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddRequest) ProtoMessage() {}

func (x *AddRequest) ProtoReflect() protoreflect.Message {
	mi := &file_example_v1_example_proto_msgTypes[0]
	if x != nil {
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
	return file_example_v1_example_proto_rawDescGZIP(), []int{0}
}

func (x *AddRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AddRequest) GetSurname() string {
	if x != nil {
		return x.Surname
	}
	return ""
}

type AddResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Result:
	//
	//	*AddResponse_User
	//	*AddResponse_Error
	Result        isAddResponse_Result `protobuf_oneof:"Result"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddResponse) Reset() {
	*x = AddResponse{}
	mi := &file_example_v1_example_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddResponse) ProtoMessage() {}

func (x *AddResponse) ProtoReflect() protoreflect.Message {
	mi := &file_example_v1_example_proto_msgTypes[1]
	if x != nil {
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
	return file_example_v1_example_proto_rawDescGZIP(), []int{1}
}

func (x *AddResponse) GetResult() isAddResponse_Result {
	if x != nil {
		return x.Result
	}
	return nil
}

func (x *AddResponse) GetUser() *User {
	if x != nil {
		if x, ok := x.Result.(*AddResponse_User); ok {
			return x.User
		}
	}
	return nil
}

func (x *AddResponse) GetError() *rpc.Error {
	if x != nil {
		if x, ok := x.Result.(*AddResponse_Error); ok {
			return x.Error
		}
	}
	return nil
}

type isAddResponse_Result interface {
	isAddResponse_Result()
}

type AddResponse_User struct {
	User *User `protobuf:"bytes,1,opt,name=user,proto3,oneof"`
}

type AddResponse_Error struct {
	Error *rpc.Error `protobuf:"bytes,2,opt,name=error,proto3,oneof"`
}

func (*AddResponse_User) isAddResponse_Result() {}

func (*AddResponse_Error) isAddResponse_Result() {}

type GetRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	mi := &file_example_v1_example_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_example_v1_example_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_example_v1_example_proto_rawDescGZIP(), []int{2}
}

func (x *GetRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Result:
	//
	//	*GetResponse_User
	//	*GetResponse_Error
	Result        isGetResponse_Result `protobuf_oneof:"Result"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetResponse) Reset() {
	*x = GetResponse{}
	mi := &file_example_v1_example_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResponse) ProtoMessage() {}

func (x *GetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_example_v1_example_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResponse.ProtoReflect.Descriptor instead.
func (*GetResponse) Descriptor() ([]byte, []int) {
	return file_example_v1_example_proto_rawDescGZIP(), []int{3}
}

func (x *GetResponse) GetResult() isGetResponse_Result {
	if x != nil {
		return x.Result
	}
	return nil
}

func (x *GetResponse) GetUser() *User {
	if x != nil {
		if x, ok := x.Result.(*GetResponse_User); ok {
			return x.User
		}
	}
	return nil
}

func (x *GetResponse) GetError() *rpc.Error {
	if x != nil {
		if x, ok := x.Result.(*GetResponse_Error); ok {
			return x.Error
		}
	}
	return nil
}

type isGetResponse_Result interface {
	isGetResponse_Result()
}

type GetResponse_User struct {
	User *User `protobuf:"bytes,1,opt,name=user,proto3,oneof"`
}

type GetResponse_Error struct {
	Error *rpc.Error `protobuf:"bytes,2,opt,name=error,proto3,oneof"`
}

func (*GetResponse_User) isGetResponse_Result() {}

func (*GetResponse_Error) isGetResponse_Result() {}

type RegistrationsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *User                  `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegistrationsResponse) Reset() {
	*x = RegistrationsResponse{}
	mi := &file_example_v1_example_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegistrationsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistrationsResponse) ProtoMessage() {}

func (x *RegistrationsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_example_v1_example_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistrationsResponse.ProtoReflect.Descriptor instead.
func (*RegistrationsResponse) Descriptor() ([]byte, []int) {
	return file_example_v1_example_proto_rawDescGZIP(), []int{4}
}

func (x *RegistrationsResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type User struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Surname       string                 `protobuf:"bytes,3,opt,name=surname,proto3" json:"surname,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *User) Reset() {
	*x = User{}
	mi := &file_example_v1_example_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_example_v1_example_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_example_v1_example_proto_rawDescGZIP(), []int{5}
}

func (x *User) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *User) GetSurname() string {
	if x != nil {
		return x.Surname
	}
	return ""
}

var File_example_v1_example_proto protoreflect.FileDescriptor

var file_example_v1_example_proto_rawDesc = string([]byte{
	0x0a, 0x18, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x65, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x73, 0x79, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65,
	0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x73, 0x79, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2f, 0x72,
	0x70, 0x63, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3a, 0x0a, 0x0a, 0x41,
	0x64, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x73, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x73, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x6d, 0x0a, 0x0b, 0x41, 0x64, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x48, 0x00, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x2c,
	0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x73, 0x79, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x48, 0x00, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x08, 0x0a, 0x06,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x1c, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x02, 0x69, 0x64, 0x22, 0x6d, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x48, 0x00, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x2c, 0x0a, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x79, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x48, 0x00, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x08, 0x0a, 0x06, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x22, 0x3d, 0x0a, 0x15, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a, 0x04,
	0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73,
	0x65, 0x72, 0x22, 0x44, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x73, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x32, 0x8d, 0x02, 0x0a, 0x0b, 0x55, 0x73, 0x65,
	0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x47, 0x0a, 0x03, 0x41, 0x64, 0x64, 0x12,
	0x16, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x0f, 0x92, 0xb5, 0x18, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x61, 0x64,
	0x64, 0x12, 0x47, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x16, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x17, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x0f, 0x92, 0xb5, 0x18, 0x0b, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x67, 0x65, 0x74, 0x12, 0x61, 0x0a, 0x0d, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x21, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x15, 0x92, 0xb5, 0x18, 0x0d, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x98, 0xb5, 0x18, 0x01, 0x1a, 0x09, 0x8a,
	0xb5, 0x18, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x42, 0xa2, 0x01, 0x0a, 0x0e, 0x63, 0x6f, 0x6d,
	0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x45, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x39, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x79, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65,
	0x74, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2d, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2d, 0x73, 0x64, 0x6b,
	0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2f, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x45, 0x58, 0x58, 0xaa, 0x02, 0x0a, 0x45,
	0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0a, 0x45, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x16, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea,
	0x02, 0x0b, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_example_v1_example_proto_rawDescOnce sync.Once
	file_example_v1_example_proto_rawDescData []byte
)

func file_example_v1_example_proto_rawDescGZIP() []byte {
	file_example_v1_example_proto_rawDescOnce.Do(func() {
		file_example_v1_example_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_example_v1_example_proto_rawDesc), len(file_example_v1_example_proto_rawDesc)))
	})
	return file_example_v1_example_proto_rawDescData
}

var file_example_v1_example_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_example_v1_example_proto_goTypes = []any{
	(*AddRequest)(nil),            // 0: example.v1.AddRequest
	(*AddResponse)(nil),           // 1: example.v1.AddResponse
	(*GetRequest)(nil),            // 2: example.v1.GetRequest
	(*GetResponse)(nil),           // 3: example.v1.GetResponse
	(*RegistrationsResponse)(nil), // 4: example.v1.RegistrationsResponse
	(*User)(nil),                  // 5: example.v1.User
	(*rpc.Error)(nil),             // 6: synternet.rpc.Error
	(*emptypb.Empty)(nil),         // 7: google.protobuf.Empty
}
var file_example_v1_example_proto_depIdxs = []int32{
	5, // 0: example.v1.AddResponse.user:type_name -> example.v1.User
	6, // 1: example.v1.AddResponse.error:type_name -> synternet.rpc.Error
	5, // 2: example.v1.GetResponse.user:type_name -> example.v1.User
	6, // 3: example.v1.GetResponse.error:type_name -> synternet.rpc.Error
	5, // 4: example.v1.RegistrationsResponse.user:type_name -> example.v1.User
	0, // 5: example.v1.UserService.Add:input_type -> example.v1.AddRequest
	2, // 6: example.v1.UserService.Get:input_type -> example.v1.GetRequest
	7, // 7: example.v1.UserService.Registrations:input_type -> google.protobuf.Empty
	1, // 8: example.v1.UserService.Add:output_type -> example.v1.AddResponse
	3, // 9: example.v1.UserService.Get:output_type -> example.v1.GetResponse
	4, // 10: example.v1.UserService.Registrations:output_type -> example.v1.RegistrationsResponse
	8, // [8:11] is the sub-list for method output_type
	5, // [5:8] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_example_v1_example_proto_init() }
func file_example_v1_example_proto_init() {
	if File_example_v1_example_proto != nil {
		return
	}
	file_example_v1_example_proto_msgTypes[1].OneofWrappers = []any{
		(*AddResponse_User)(nil),
		(*AddResponse_Error)(nil),
	}
	file_example_v1_example_proto_msgTypes[3].OneofWrappers = []any{
		(*GetResponse_User)(nil),
		(*GetResponse_Error)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_example_v1_example_proto_rawDesc), len(file_example_v1_example_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_example_v1_example_proto_goTypes,
		DependencyIndexes: file_example_v1_example_proto_depIdxs,
		MessageInfos:      file_example_v1_example_proto_msgTypes,
	}.Build()
	File_example_v1_example_proto = out.File
	file_example_v1_example_proto_goTypes = nil
	file_example_v1_example_proto_depIdxs = nil
}
