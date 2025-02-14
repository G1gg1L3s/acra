// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: encryptedTypes.proto

package common

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

// EncryptedType defines types for encrypted data.
type EncryptedType int32

const (
	EncryptedType_Unknown EncryptedType = 0
	EncryptedType_Int32   EncryptedType = 1
	EncryptedType_Int64   EncryptedType = 2
	EncryptedType_String  EncryptedType = 3
	EncryptedType_Bytes   EncryptedType = 4
)

// Enum value maps for EncryptedType.
var (
	EncryptedType_name = map[int32]string{
		0: "Unknown",
		1: "Int32",
		2: "Int64",
		3: "String",
		4: "Bytes",
	}
	EncryptedType_value = map[string]int32{
		"Unknown": 0,
		"Int32":   1,
		"Int64":   2,
		"String":  3,
		"Bytes":   4,
	}
)

func (x EncryptedType) Enum() *EncryptedType {
	p := new(EncryptedType)
	*p = x
	return p
}

func (x EncryptedType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EncryptedType) Descriptor() protoreflect.EnumDescriptor {
	return file_encryptedTypes_proto_enumTypes[0].Descriptor()
}

func (EncryptedType) Type() protoreflect.EnumType {
	return &file_encryptedTypes_proto_enumTypes[0]
}

func (x EncryptedType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EncryptedType.Descriptor instead.
func (EncryptedType) EnumDescriptor() ([]byte, []int) {
	return file_encryptedTypes_proto_rawDescGZIP(), []int{0}
}

// EncryptedValue keeps serialized encrypted value.
type EncryptedValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []byte        `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Type  EncryptedType `protobuf:"varint,2,opt,name=type,proto3,enum=github.com.cossacklabs.acra.encryptor.config.common.EncryptedType" json:"type,omitempty"`
}

func (x *EncryptedValue) Reset() {
	*x = EncryptedValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_encryptedTypes_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EncryptedValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EncryptedValue) ProtoMessage() {}

func (x *EncryptedValue) ProtoReflect() protoreflect.Message {
	mi := &file_encryptedTypes_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EncryptedValue.ProtoReflect.Descriptor instead.
func (*EncryptedValue) Descriptor() ([]byte, []int) {
	return file_encryptedTypes_proto_rawDescGZIP(), []int{0}
}

func (x *EncryptedValue) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *EncryptedValue) GetType() EncryptedType {
	if x != nil {
		return x.Type
	}
	return EncryptedType_Unknown
}

var File_encryptedTypes_proto protoreflect.FileDescriptor

var file_encryptedTypes_proto_rawDesc = []byte{
	0x0a, 0x14, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x54, 0x79, 0x70, 0x65, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x73, 0x73, 0x61, 0x63, 0x6b, 0x6c, 0x61, 0x62, 0x73, 0x2e, 0x61,
	0x63, 0x72, 0x61, 0x2e, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x22, 0x7e, 0x0a, 0x0e, 0x45,
	0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x56, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x42, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x63,
	0x6f, 0x73, 0x73, 0x61, 0x63, 0x6b, 0x6c, 0x61, 0x62, 0x73, 0x2e, 0x61, 0x63, 0x72, 0x61, 0x2e,
	0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65,
	0x64, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x2a, 0x49, 0x0a, 0x0d, 0x45,
	0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07,
	0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x49, 0x6e, 0x74,
	0x33, 0x32, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x10, 0x02, 0x12,
	0x0a, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x10, 0x03, 0x12, 0x09, 0x0a, 0x05, 0x42,
	0x79, 0x74, 0x65, 0x73, 0x10, 0x04, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x73, 0x73, 0x61, 0x63, 0x6b, 0x6c, 0x61, 0x62, 0x73,
	0x2f, 0x61, 0x63, 0x72, 0x61, 0x2f, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x72, 0x2f,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_encryptedTypes_proto_rawDescOnce sync.Once
	file_encryptedTypes_proto_rawDescData = file_encryptedTypes_proto_rawDesc
)

func file_encryptedTypes_proto_rawDescGZIP() []byte {
	file_encryptedTypes_proto_rawDescOnce.Do(func() {
		file_encryptedTypes_proto_rawDescData = protoimpl.X.CompressGZIP(file_encryptedTypes_proto_rawDescData)
	})
	return file_encryptedTypes_proto_rawDescData
}

var file_encryptedTypes_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_encryptedTypes_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_encryptedTypes_proto_goTypes = []interface{}{
	(EncryptedType)(0),     // 0: github.com.cossacklabs.acra.encryptor.config.common.EncryptedType
	(*EncryptedValue)(nil), // 1: github.com.cossacklabs.acra.encryptor.config.common.EncryptedValue
}
var file_encryptedTypes_proto_depIdxs = []int32{
	0, // 0: github.com.cossacklabs.acra.encryptor.config.common.EncryptedValue.type:type_name -> github.com.cossacklabs.acra.encryptor.config.common.EncryptedType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_encryptedTypes_proto_init() }
func file_encryptedTypes_proto_init() {
	if File_encryptedTypes_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_encryptedTypes_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EncryptedValue); i {
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
			RawDescriptor: file_encryptedTypes_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_encryptedTypes_proto_goTypes,
		DependencyIndexes: file_encryptedTypes_proto_depIdxs,
		EnumInfos:         file_encryptedTypes_proto_enumTypes,
		MessageInfos:      file_encryptedTypes_proto_msgTypes,
	}.Build()
	File_encryptedTypes_proto = out.File
	file_encryptedTypes_proto_rawDesc = nil
	file_encryptedTypes_proto_goTypes = nil
	file_encryptedTypes_proto_depIdxs = nil
}
