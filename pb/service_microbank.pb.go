// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: service_microbank.proto

package pb

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_service_microbank_proto protoreflect.FileDescriptor

var file_service_microbank_proto_rawDesc = []byte{
	0x0a, 0x17, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x62,
	0x61, 0x6e, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70,
	0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x72, 0x70, 0x63,
	0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x14, 0x72, 0x70, 0x63, 0x5f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x75, 0x73,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x72, 0x70, 0x63, 0x5f, 0x73, 0x65,
	0x6e, 0x64, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x32, 0xea, 0x03, 0x0a, 0x09, 0x4d, 0x69, 0x63, 0x72, 0x6f, 0x42, 0x61, 0x6e, 0x6b, 0x12,
	0x8e, 0x01, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x15,
	0x2e, 0x70, 0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x51, 0x92,
	0x41, 0x34, 0x12, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x20, 0x6e, 0x65, 0x77, 0x20, 0x75,
	0x73, 0x65, 0x72, 0x1a, 0x21, 0x55, 0x73, 0x65, 0x20, 0x74, 0x68, 0x69, 0x73, 0x20, 0x41, 0x50,
	0x49, 0x20, 0x74, 0x6f, 0x20, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x20, 0x61, 0x20, 0x6e, 0x65,
	0x77, 0x20, 0x75, 0x73, 0x65, 0x72, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x3a, 0x01, 0x2a, 0x22,
	0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x12, 0xa3, 0x01, 0x0a, 0x09, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x14,
	0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x69, 0x92, 0x41, 0x4d,
	0x12, 0x0a, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x20, 0x75, 0x73, 0x65, 0x72, 0x1a, 0x3f, 0x55, 0x73,
	0x65, 0x20, 0x74, 0x68, 0x69, 0x73, 0x20, 0x41, 0x50, 0x49, 0x20, 0x74, 0x6f, 0x20, 0x6c, 0x6f,
	0x67, 0x69, 0x6e, 0x20, 0x75, 0x73, 0x65, 0x72, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x67, 0x65, 0x74,
	0x20, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x20, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x20, 0x26, 0x20,
	0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x20, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x13, 0x3a, 0x01, 0x2a, 0x22, 0x0e, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72,
	0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0xa5, 0x01, 0x0a, 0x0c, 0x53, 0x65, 0x6e, 0x64, 0x50,
	0x61, 0x73, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x17, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x6e,
	0x64, 0x50, 0x61, 0x73, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x18, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x50, 0x61, 0x73, 0x73, 0x43, 0x6f,
	0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x62, 0x92, 0x41, 0x3e, 0x12,
	0x0d, 0x53, 0x65, 0x6e, 0x64, 0x20, 0x70, 0x61, 0x73, 0x73, 0x63, 0x6f, 0x64, 0x65, 0x1a, 0x2d,
	0x55, 0x73, 0x65, 0x20, 0x74, 0x68, 0x69, 0x73, 0x20, 0x41, 0x50, 0x49, 0x20, 0x74, 0x6f, 0x20,
	0x73, 0x65, 0x6e, 0x64, 0x20, 0x70, 0x61, 0x73, 0x73, 0x63, 0x6f, 0x64, 0x65, 0x20, 0x74, 0x6f,
	0x20, 0x75, 0x73, 0x65, 0x72, 0x27, 0x73, 0x20, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x1b, 0x3a, 0x01, 0x2a, 0x22, 0x16, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72,
	0x2f, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x63, 0x6f, 0x64, 0x65, 0x42, 0x83,
	0x01, 0x92, 0x41, 0x59, 0x12, 0x57, 0x0a, 0x0e, 0x4d, 0x69, 0x63, 0x72, 0x6f, 0x20, 0x42, 0x61,
	0x6e, 0x6b, 0x20, 0x41, 0x50, 0x49, 0x22, 0x40, 0x0a, 0x0d, 0x43, 0x6f, 0x64, 0x65, 0x53, 0x69,
	0x6e, 0x67, 0x65, 0x72, 0x47, 0x6e, 0x43, 0x12, 0x20, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f,
	0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x43, 0x6f, 0x64, 0x65,
	0x53, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x47, 0x6e, 0x43, 0x1a, 0x0d, 0x75, 0x6e, 0x6b, 0x6e, 0x6f,
	0x77, 0x40, 0x71, 0x71, 0x2e, 0x63, 0x6f, 0x6d, 0x32, 0x03, 0x31, 0x2e, 0x32, 0x5a, 0x25, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x43, 0x6f, 0x64, 0x65, 0x53, 0x69,
	0x6e, 0x67, 0x65, 0x72, 0x47, 0x6e, 0x43, 0x2f, 0x4d, 0x69, 0x63, 0x72, 0x6f, 0x42, 0x61, 0x6e,
	0x6b, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_service_microbank_proto_goTypes = []any{
	(*CreateUserRequest)(nil),    // 0: pb.CreateUserRequest
	(*LoginUserRequest)(nil),     // 1: pb.LoginUserRequest
	(*SendPassCodeRequest)(nil),  // 2: pb.SendPassCodeRequest
	(*CreateUserResponse)(nil),   // 3: pb.CreateUserResponse
	(*LoginUserResponse)(nil),    // 4: pb.LoginUserResponse
	(*SendPassCodeResponse)(nil), // 5: pb.SendPassCodeResponse
}
var file_service_microbank_proto_depIdxs = []int32{
	0, // 0: pb.MicroBank.CreateUser:input_type -> pb.CreateUserRequest
	1, // 1: pb.MicroBank.LoginUser:input_type -> pb.LoginUserRequest
	2, // 2: pb.MicroBank.SendPassCode:input_type -> pb.SendPassCodeRequest
	3, // 3: pb.MicroBank.CreateUser:output_type -> pb.CreateUserResponse
	4, // 4: pb.MicroBank.LoginUser:output_type -> pb.LoginUserResponse
	5, // 5: pb.MicroBank.SendPassCode:output_type -> pb.SendPassCodeResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_service_microbank_proto_init() }
func file_service_microbank_proto_init() {
	if File_service_microbank_proto != nil {
		return
	}
	file_rpc_create_user_proto_init()
	file_rpc_login_user_proto_init()
	file_rpc_send_passcode_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_service_microbank_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_microbank_proto_goTypes,
		DependencyIndexes: file_service_microbank_proto_depIdxs,
	}.Build()
	File_service_microbank_proto = out.File
	file_service_microbank_proto_rawDesc = nil
	file_service_microbank_proto_goTypes = nil
	file_service_microbank_proto_depIdxs = nil
}
