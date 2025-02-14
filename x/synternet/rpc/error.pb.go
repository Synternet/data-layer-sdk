// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: synternet/rpc/error.proto

package rpc

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

// Error message can be used to decode error sent in a message
type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,400,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_synternet_rpc_error_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_synternet_rpc_error_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_synternet_rpc_error_proto_rawDescGZIP(), []int{0}
}

func (x *Error) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_synternet_rpc_error_proto protoreflect.FileDescriptor

var file_synternet_rpc_error_proto_rawDesc = []byte{
	0x0a, 0x19, 0x73, 0x79, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x73, 0x79, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x22, 0x1e, 0x0a, 0x05, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x12, 0x15, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x90, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x42, 0xa9, 0x01, 0x0a, 0x11, 0x63,
	0x6f, 0x6d, 0x2e, 0x73, 0x79, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e, 0x72, 0x70, 0x63,
	0x42, 0x0a, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x33,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x79, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x65, 0x74, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2d, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2d,
	0x73, 0x64, 0x6b, 0x2f, 0x78, 0x2f, 0x73, 0x79, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2f,
	0x72, 0x70, 0x63, 0xa2, 0x02, 0x03, 0x53, 0x52, 0x58, 0xaa, 0x02, 0x0d, 0x53, 0x79, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e, 0x52, 0x70, 0x63, 0xca, 0x02, 0x0d, 0x53, 0x79, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x65, 0x74, 0x5c, 0x52, 0x70, 0x63, 0xe2, 0x02, 0x19, 0x53, 0x79, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x65, 0x74, 0x5c, 0x52, 0x70, 0x63, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0e, 0x53, 0x79, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65,
	0x74, 0x3a, 0x3a, 0x52, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_synternet_rpc_error_proto_rawDescOnce sync.Once
	file_synternet_rpc_error_proto_rawDescData = file_synternet_rpc_error_proto_rawDesc
)

func file_synternet_rpc_error_proto_rawDescGZIP() []byte {
	file_synternet_rpc_error_proto_rawDescOnce.Do(func() {
		file_synternet_rpc_error_proto_rawDescData = protoimpl.X.CompressGZIP(file_synternet_rpc_error_proto_rawDescData)
	})
	return file_synternet_rpc_error_proto_rawDescData
}

var file_synternet_rpc_error_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_synternet_rpc_error_proto_goTypes = []interface{}{
	(*Error)(nil), // 0: synternet.rpc.Error
}
var file_synternet_rpc_error_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_synternet_rpc_error_proto_init() }
func file_synternet_rpc_error_proto_init() {
	if File_synternet_rpc_error_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_synternet_rpc_error_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
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
			RawDescriptor: file_synternet_rpc_error_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_synternet_rpc_error_proto_goTypes,
		DependencyIndexes: file_synternet_rpc_error_proto_depIdxs,
		MessageInfos:      file_synternet_rpc_error_proto_msgTypes,
	}.Build()
	File_synternet_rpc_error_proto = out.File
	file_synternet_rpc_error_proto_rawDesc = nil
	file_synternet_rpc_error_proto_goTypes = nil
	file_synternet_rpc_error_proto_depIdxs = nil
}
