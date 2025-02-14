// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: types/rpc/options.proto

package rpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var file_types_rpc_options_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.ServiceOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         50001,
		Name:          "types.rpc.subject_prefix",
		Tag:           "bytes,50001,opt,name=subject_prefix",
		Filename:      "types/rpc/options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         50002,
		Name:          "types.rpc.subject_suffix",
		Tag:           "bytes,50002,opt,name=subject_suffix",
		Filename:      "types/rpc/options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50003,
		Name:          "types.rpc.skip_inputs",
		Tag:           "varint,50003,opt,name=skip_inputs",
		Filename:      "types/rpc/options.proto",
	},
}

// Extension fields to descriptorpb.ServiceOptions.
var (
	// Optional prefix for all methods in a service
	//
	// optional string subject_prefix = 50001;
	E_SubjectPrefix = &file_types_rpc_options_proto_extTypes[0]
)

// Extension fields to descriptorpb.MethodOptions.
var (
	// Optional suffix to override default method name in a subject
	//
	// optional string subject_suffix = 50002;
	E_SubjectSuffix = &file_types_rpc_options_proto_extTypes[1]
	// Optional flag to skip subscribing to the input stream.
	//
	// optional bool skip_inputs = 50003;
	E_SkipInputs = &file_types_rpc_options_proto_extTypes[2]
)

var File_types_rpc_options_proto protoreflect.FileDescriptor

var file_types_rpc_options_proto_rawDesc = []byte{
	0x0a, 0x17, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2e, 0x72, 0x70, 0x63, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x48, 0x0a, 0x0e, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x5f, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd1, 0x86, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78,
	0x3a, 0x47, 0x0a, 0x0e, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x73, 0x75, 0x66, 0x66,
	0x69, 0x78, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0xd2, 0x86, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x75, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x53, 0x75, 0x66, 0x66, 0x69, 0x78, 0x3a, 0x41, 0x0a, 0x0b, 0x73, 0x6b, 0x69,
	0x70, 0x5f, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x73, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd3, 0x86, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0a, 0x73, 0x6b, 0x69, 0x70, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x73, 0x42, 0x91, 0x01, 0x0a,
	0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x72, 0x70, 0x63, 0x42, 0x0c,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2d,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x79, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x65, 0x74, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2d, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2d,
	0x73, 0x64, 0x6b, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x72, 0x70, 0x63, 0xa2, 0x02, 0x03,
	0x54, 0x52, 0x58, 0xaa, 0x02, 0x09, 0x54, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x52, 0x70, 0x63, 0xca,
	0x02, 0x09, 0x54, 0x79, 0x70, 0x65, 0x73, 0x5c, 0x52, 0x70, 0x63, 0xe2, 0x02, 0x15, 0x54, 0x79,
	0x70, 0x65, 0x73, 0x5c, 0x52, 0x70, 0x63, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0xea, 0x02, 0x0a, 0x54, 0x79, 0x70, 0x65, 0x73, 0x3a, 0x3a, 0x52, 0x70, 0x63,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_types_rpc_options_proto_goTypes = []interface{}{
	(*descriptorpb.ServiceOptions)(nil), // 0: google.protobuf.ServiceOptions
	(*descriptorpb.MethodOptions)(nil),  // 1: google.protobuf.MethodOptions
}
var file_types_rpc_options_proto_depIdxs = []int32{
	0, // 0: types.rpc.subject_prefix:extendee -> google.protobuf.ServiceOptions
	1, // 1: types.rpc.subject_suffix:extendee -> google.protobuf.MethodOptions
	1, // 2: types.rpc.skip_inputs:extendee -> google.protobuf.MethodOptions
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	0, // [0:3] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_types_rpc_options_proto_init() }
func file_types_rpc_options_proto_init() {
	if File_types_rpc_options_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_types_rpc_options_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 3,
			NumServices:   0,
		},
		GoTypes:           file_types_rpc_options_proto_goTypes,
		DependencyIndexes: file_types_rpc_options_proto_depIdxs,
		ExtensionInfos:    file_types_rpc_options_proto_extTypes,
	}.Build()
	File_types_rpc_options_proto = out.File
	file_types_rpc_options_proto_rawDesc = nil
	file_types_rpc_options_proto_goTypes = nil
	file_types_rpc_options_proto_depIdxs = nil
}
