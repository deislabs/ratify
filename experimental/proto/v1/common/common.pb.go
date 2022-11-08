// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: common.proto

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

// Descriptor holds various properties of an artifact.
//
// (oci.Descriptor)
//
// The only strictly defined property is `rawPath`. Beyond that, it is designed to be flexible: only what is known is included.
type Descriptor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The path or address that can be used to retrieve the artifact
	RawPath string `protobuf:"bytes,1,opt,name=rawPath,proto3" json:"rawPath,omitempty"`
	// Collection of key/value pairs containing relevant known metadata
	// Examples would include tag or digest.
	Attributes []*Descriptor_Attributes `protobuf:"bytes,2,rep,name=attributes,proto3" json:"attributes,omitempty"`
}

func (x *Descriptor) Reset() {
	*x = Descriptor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Descriptor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Descriptor) ProtoMessage() {}

func (x *Descriptor) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Descriptor.ProtoReflect.Descriptor instead.
func (*Descriptor) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{0}
}

func (x *Descriptor) GetRawPath() string {
	if x != nil {
		return x.RawPath
	}
	return ""
}

func (x *Descriptor) GetAttributes() []*Descriptor_Attributes {
	if x != nil {
		return x.Attributes
	}
	return nil
}

// Referrer holds properties on a given Artifact.
//
// (Ratify.ocispecs.ReferenceDescriptor)
type Referrer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The type of the artifact.
	ArtifactType string `protobuf:"bytes,1,opt,name=artifactType,proto3" json:"artifactType,omitempty"`
	// All known properties (parsed or retrieved) for the artifact.
	Descriptor_ *Descriptor `protobuf:"bytes,2,opt,name=descriptor,proto3" json:"descriptor,omitempty"`
}

func (x *Referrer) Reset() {
	*x = Referrer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Referrer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Referrer) ProtoMessage() {}

func (x *Referrer) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Referrer.ProtoReflect.Descriptor instead.
func (*Referrer) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{1}
}

func (x *Referrer) GetArtifactType() string {
	if x != nil {
		return x.ArtifactType
	}
	return ""
}

func (x *Referrer) GetDescriptor_() *Descriptor {
	if x != nil {
		return x.Descriptor_
	}
	return nil
}

// Manifest includes all relevant properties and blobs for a reference artifact.
//
// (Ratify.ocispecs.ReferenceManifest)
type Manifest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The properties of the referrer.
	Referrer *Referrer `protobuf:"bytes,1,opt,name=referrer,proto3" json:"referrer,omitempty"`
	// The collection of layers which together form the content of the referrer.
	Blobs *Manifest_Blobs `protobuf:"bytes,2,opt,name=blobs,proto3" json:"blobs,omitempty"`
	// The parents (subjects) for which this referrer describes and supplements.
	Subject *Manifest_Subjects `protobuf:"bytes,3,opt,name=subject,proto3" json:"subject,omitempty"`
}

func (x *Manifest) Reset() {
	*x = Manifest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Manifest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Manifest) ProtoMessage() {}

func (x *Manifest) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Manifest.ProtoReflect.Descriptor instead.
func (*Manifest) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{2}
}

func (x *Manifest) GetReferrer() *Referrer {
	if x != nil {
		return x.Referrer
	}
	return nil
}

func (x *Manifest) GetBlobs() *Manifest_Blobs {
	if x != nil {
		return x.Blobs
	}
	return nil
}

func (x *Manifest) GetSubject() *Manifest_Subjects {
	if x != nil {
		return x.Subject
	}
	return nil
}

type Descriptor_Attributes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values map[string]string `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Descriptor_Attributes) Reset() {
	*x = Descriptor_Attributes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Descriptor_Attributes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Descriptor_Attributes) ProtoMessage() {}

func (x *Descriptor_Attributes) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Descriptor_Attributes.ProtoReflect.Descriptor instead.
func (*Descriptor_Attributes) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Descriptor_Attributes) GetValues() map[string]string {
	if x != nil {
		return x.Values
	}
	return nil
}

type Manifest_Blobs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Descriptor_ []*Descriptor `protobuf:"bytes,1,rep,name=descriptor,proto3" json:"descriptor,omitempty"`
}

func (x *Manifest_Blobs) Reset() {
	*x = Manifest_Blobs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Manifest_Blobs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Manifest_Blobs) ProtoMessage() {}

func (x *Manifest_Blobs) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Manifest_Blobs.ProtoReflect.Descriptor instead.
func (*Manifest_Blobs) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{2, 0}
}

func (x *Manifest_Blobs) GetDescriptor_() []*Descriptor {
	if x != nil {
		return x.Descriptor_
	}
	return nil
}

type Manifest_Subjects struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Descriptor_ []*Descriptor `protobuf:"bytes,1,rep,name=descriptor,proto3" json:"descriptor,omitempty"`
}

func (x *Manifest_Subjects) Reset() {
	*x = Manifest_Subjects{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Manifest_Subjects) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Manifest_Subjects) ProtoMessage() {}

func (x *Manifest_Subjects) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Manifest_Subjects.ProtoReflect.Descriptor instead.
func (*Manifest_Subjects) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{2, 1}
}

func (x *Manifest_Subjects) GetDescriptor_() []*Descriptor {
	if x != nil {
		return x.Descriptor_
	}
	return nil
}

var File_common_proto protoreflect.FileDescriptor

var file_common_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x22, 0xf2, 0x01, 0x0a, 0x0a, 0x44, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x61, 0x77, 0x50, 0x61, 0x74, 0x68,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x72, 0x61, 0x77, 0x50, 0x61, 0x74, 0x68, 0x12,
	0x3d, 0x0a, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x44, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x73, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x1a, 0x8a,
	0x01, 0x0a, 0x0a, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x41, 0x0a,
	0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f,
	0x72, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x2e, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x1a, 0x39, 0x0a, 0x0b, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x62, 0x0a, 0x08, 0x52,
	0x65, 0x66, 0x65, 0x72, 0x72, 0x65, 0x72, 0x12, 0x22, 0x0a, 0x0c, 0x61, 0x72, 0x74, 0x69, 0x66,
	0x61, 0x63, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x61,
	0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x32, 0x0a, 0x0a, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x6f, 0x72, 0x52, 0x0a, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x22,
	0x98, 0x02, 0x0a, 0x08, 0x4d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x08,
	0x72, 0x65, 0x66, 0x65, 0x72, 0x72, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x72, 0x65, 0x72,
	0x52, 0x08, 0x72, 0x65, 0x66, 0x65, 0x72, 0x72, 0x65, 0x72, 0x12, 0x2c, 0x0a, 0x05, 0x62, 0x6c,
	0x6f, 0x62, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x4d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x2e, 0x42, 0x6c, 0x6f, 0x62,
	0x73, 0x52, 0x05, 0x62, 0x6c, 0x6f, 0x62, 0x73, 0x12, 0x33, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x4d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x2e, 0x53, 0x75, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x73, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x1a, 0x3b, 0x0a,
	0x05, 0x42, 0x6c, 0x6f, 0x62, 0x73, 0x12, 0x32, 0x0a, 0x0a, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x52, 0x0a,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x1a, 0x3e, 0x0a, 0x08, 0x53, 0x75,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x12, 0x32, 0x0a, 0x0a, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x52, 0x0a,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x65, 0x69, 0x73, 0x6c, 0x61, 0x62,
	0x73, 0x2f, 0x72, 0x61, 0x74, 0x69, 0x66, 0x79, 0x2f, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d,
	0x65, 0x6e, 0x74, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_proto_rawDescOnce sync.Once
	file_common_proto_rawDescData = file_common_proto_rawDesc
)

func file_common_proto_rawDescGZIP() []byte {
	file_common_proto_rawDescOnce.Do(func() {
		file_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_proto_rawDescData)
	})
	return file_common_proto_rawDescData
}

var file_common_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_common_proto_goTypes = []interface{}{
	(*Descriptor)(nil),            // 0: common.Descriptor
	(*Referrer)(nil),              // 1: common.Referrer
	(*Manifest)(nil),              // 2: common.Manifest
	(*Descriptor_Attributes)(nil), // 3: common.Descriptor.Attributes
	nil,                           // 4: common.Descriptor.Attributes.ValuesEntry
	(*Manifest_Blobs)(nil),        // 5: common.Manifest.Blobs
	(*Manifest_Subjects)(nil),     // 6: common.Manifest.Subjects
}
var file_common_proto_depIdxs = []int32{
	3, // 0: common.Descriptor.attributes:type_name -> common.Descriptor.Attributes
	0, // 1: common.Referrer.descriptor:type_name -> common.Descriptor
	1, // 2: common.Manifest.referrer:type_name -> common.Referrer
	5, // 3: common.Manifest.blobs:type_name -> common.Manifest.Blobs
	6, // 4: common.Manifest.subject:type_name -> common.Manifest.Subjects
	4, // 5: common.Descriptor.Attributes.values:type_name -> common.Descriptor.Attributes.ValuesEntry
	0, // 6: common.Manifest.Blobs.descriptor:type_name -> common.Descriptor
	0, // 7: common.Manifest.Subjects.descriptor:type_name -> common.Descriptor
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_common_proto_init() }
func file_common_proto_init() {
	if File_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Descriptor); i {
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
		file_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Referrer); i {
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
		file_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Manifest); i {
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
		file_common_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Descriptor_Attributes); i {
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
		file_common_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Manifest_Blobs); i {
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
		file_common_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Manifest_Subjects); i {
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
			RawDescriptor: file_common_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_proto_goTypes,
		DependencyIndexes: file_common_proto_depIdxs,
		MessageInfos:      file_common_proto_msgTypes,
	}.Build()
	File_common_proto = out.File
	file_common_proto_rawDesc = nil
	file_common_proto_goTypes = nil
	file_common_proto_depIdxs = nil
}
