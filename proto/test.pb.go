// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.0
// source: test.proto

package __

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

type Test1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	I *int32  `protobuf:"varint,1,req,name=I" json:"I,omitempty"`
	S *string `protobuf:"bytes,2,req,name=S" json:"S,omitempty"`
	B *bool   `protobuf:"varint,3,req,name=B" json:"B,omitempty"`
}

func (x *Test1) Reset() {
	*x = Test1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Test1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Test1) ProtoMessage() {}

func (x *Test1) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Test1.ProtoReflect.Descriptor instead.
func (*Test1) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{0}
}

func (x *Test1) GetI() int32 {
	if x != nil && x.I != nil {
		return *x.I
	}
	return 0
}

func (x *Test1) GetS() string {
	if x != nil && x.S != nil {
		return *x.S
	}
	return ""
}

func (x *Test1) GetB() bool {
	if x != nil && x.B != nil {
		return *x.B
	}
	return false
}

type Test2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	R []*Test1 `protobuf:"bytes,1,rep,name=R" json:"R,omitempty"`
}

func (x *Test2) Reset() {
	*x = Test2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Test2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Test2) ProtoMessage() {}

func (x *Test2) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Test2.ProtoReflect.Descriptor instead.
func (*Test2) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{1}
}

func (x *Test2) GetR() []*Test1 {
	if x != nil {
		return x.R
	}
	return nil
}

var File_test_proto protoreflect.FileDescriptor

var file_test_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x74, 0x65,
	0x73, 0x74, 0x22, 0x31, 0x0a, 0x05, 0x54, 0x65, 0x73, 0x74, 0x31, 0x12, 0x0c, 0x0a, 0x01, 0x49,
	0x18, 0x01, 0x20, 0x02, 0x28, 0x05, 0x52, 0x01, 0x49, 0x12, 0x0c, 0x0a, 0x01, 0x53, 0x18, 0x02,
	0x20, 0x02, 0x28, 0x09, 0x52, 0x01, 0x53, 0x12, 0x0c, 0x0a, 0x01, 0x42, 0x18, 0x03, 0x20, 0x02,
	0x28, 0x08, 0x52, 0x01, 0x42, 0x22, 0x22, 0x0a, 0x05, 0x54, 0x65, 0x73, 0x74, 0x32, 0x12, 0x19,
	0x0a, 0x01, 0x52, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x74, 0x65, 0x73, 0x74,
	0x2e, 0x54, 0x65, 0x73, 0x74, 0x31, 0x52, 0x01, 0x52, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f,
}

var (
	file_test_proto_rawDescOnce sync.Once
	file_test_proto_rawDescData = file_test_proto_rawDesc
)

func file_test_proto_rawDescGZIP() []byte {
	file_test_proto_rawDescOnce.Do(func() {
		file_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_test_proto_rawDescData)
	})
	return file_test_proto_rawDescData
}

var file_test_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_test_proto_goTypes = []interface{}{
	(*Test1)(nil), // 0: test.Test1
	(*Test2)(nil), // 1: test.Test2
}
var file_test_proto_depIdxs = []int32{
	0, // 0: test.Test2.R:type_name -> test.Test1
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_test_proto_init() }
func file_test_proto_init() {
	if File_test_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_test_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Test1); i {
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
		file_test_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Test2); i {
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
			RawDescriptor: file_test_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_test_proto_goTypes,
		DependencyIndexes: file_test_proto_depIdxs,
		MessageInfos:      file_test_proto_msgTypes,
	}.Build()
	File_test_proto = out.File
	file_test_proto_rawDesc = nil
	file_test_proto_goTypes = nil
	file_test_proto_depIdxs = nil
}
