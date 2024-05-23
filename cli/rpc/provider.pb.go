// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: provider.proto

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

type PingReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cert  string `protobuf:"bytes,1,opt,name=cert,proto3" json:"cert,omitempty"`
	Quote []byte `protobuf:"bytes,2,opt,name=quote,proto3" json:"quote,omitempty"`
}

func (x *PingReq) Reset() {
	*x = PingReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provider_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingReq) ProtoMessage() {}

func (x *PingReq) ProtoReflect() protoreflect.Message {
	mi := &file_provider_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingReq.ProtoReflect.Descriptor instead.
func (*PingReq) Descriptor() ([]byte, []int) {
	return file_provider_proto_rawDescGZIP(), []int{0}
}

func (x *PingReq) GetCert() string {
	if x != nil {
		return x.Cert
	}
	return ""
}

func (x *PingReq) GetQuote() []byte {
	if x != nil {
		return x.Quote
	}
	return nil
}

type PingResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *PingResp) Reset() {
	*x = PingResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provider_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResp) ProtoMessage() {}

func (x *PingResp) ProtoReflect() protoreflect.Message {
	mi := &file_provider_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingResp.ProtoReflect.Descriptor instead.
func (*PingResp) Descriptor() ([]byte, []int) {
	return file_provider_proto_rawDescGZIP(), []int{1}
}

func (x *PingResp) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

var File_provider_proto protoreflect.FileDescriptor

var file_provider_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0b, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x50, 0x43, 0x22, 0x33, 0x0a,
	0x07, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x65, 0x72, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x65, 0x72, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x71, 0x75, 0x6f, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x71, 0x75, 0x6f,
	0x74, 0x65, 0x22, 0x1a, 0x0a, 0x08, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x12, 0x0e,
	0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x32, 0x3f,
	0x0a, 0x08, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x33, 0x0a, 0x04, 0x50, 0x69,
	0x6e, 0x67, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x50, 0x43,
	0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x52, 0x50, 0x43, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x42,
	0x06, 0x5a, 0x04, 0x72, 0x70, 0x63, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_provider_proto_rawDescOnce sync.Once
	file_provider_proto_rawDescData = file_provider_proto_rawDesc
)

func file_provider_proto_rawDescGZIP() []byte {
	file_provider_proto_rawDescOnce.Do(func() {
		file_provider_proto_rawDescData = protoimpl.X.CompressGZIP(file_provider_proto_rawDescData)
	})
	return file_provider_proto_rawDescData
}

var file_provider_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_provider_proto_goTypes = []interface{}{
	(*PingReq)(nil),  // 0: providerRPC.PingReq
	(*PingResp)(nil), // 1: providerRPC.PingResp
}
var file_provider_proto_depIdxs = []int32{
	0, // 0: providerRPC.Provider.Ping:input_type -> providerRPC.PingReq
	1, // 1: providerRPC.Provider.Ping:output_type -> providerRPC.PingResp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_provider_proto_init() }
func file_provider_proto_init() {
	if File_provider_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_provider_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingReq); i {
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
		file_provider_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingResp); i {
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
			RawDescriptor: file_provider_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_provider_proto_goTypes,
		DependencyIndexes: file_provider_proto_depIdxs,
		MessageInfos:      file_provider_proto_msgTypes,
	}.Build()
	File_provider_proto = out.File
	file_provider_proto_rawDesc = nil
	file_provider_proto_goTypes = nil
	file_provider_proto_depIdxs = nil
}
