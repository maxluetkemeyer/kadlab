// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.12.4
// source: proto/internal.proto

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

type StoreResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=Success,proto3" json:"Success,omitempty"`
}

func (x *StoreResult) Reset() {
	*x = StoreResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_kademlia_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreResult) ProtoMessage() {}

func (x *StoreResult) ProtoReflect() protoreflect.Message {
	mi := &file_proto_kademlia_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreResult.ProtoReflect.Descriptor instead.
func (*StoreResult) Descriptor() ([]byte, []int) {
	return file_proto_kademlia_proto_rawDescGZIP(), []int{0}
}

func (x *StoreResult) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type Content struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   *KademliaID `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value string      `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *Content) Reset() {
	*x = Content{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_kademlia_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Content) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Content) ProtoMessage() {}

func (x *Content) ProtoReflect() protoreflect.Message {
	mi := &file_proto_kademlia_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Content.ProtoReflect.Descriptor instead.
func (*Content) Descriptor() ([]byte, []int) {
	return file_proto_kademlia_proto_rawDescGZIP(), []int{1}
}

func (x *Content) GetKey() *KademliaID {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *Content) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type KademliaID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []byte `protobuf:"bytes,1,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *KademliaID) Reset() {
	*x = KademliaID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_kademlia_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KademliaID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KademliaID) ProtoMessage() {}

func (x *KademliaID) ProtoReflect() protoreflect.Message {
	mi := &file_proto_kademlia_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KademliaID.ProtoReflect.Descriptor instead.
func (*KademliaID) Descriptor() ([]byte, []int) {
	return file_proto_kademlia_proto_rawDescGZIP(), []int{2}
}

func (x *KademliaID) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID         *KademliaID `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	IPWithPort string      `protobuf:"bytes,2,opt,name=IPWithPort,proto3" json:"IPWithPort,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_kademlia_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_proto_kademlia_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_proto_kademlia_proto_rawDescGZIP(), []int{3}
}

func (x *Node) GetID() *KademliaID {
	if x != nil {
		return x.ID
	}
	return nil
}

func (x *Node) GetIPWithPort() string {
	if x != nil {
		return x.IPWithPort
	}
	return ""
}

type Nodes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Node []*Node `protobuf:"bytes,1,rep,name=Node,proto3" json:"Node,omitempty"`
}

func (x *Nodes) Reset() {
	*x = Nodes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_kademlia_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Nodes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Nodes) ProtoMessage() {}

func (x *Nodes) ProtoReflect() protoreflect.Message {
	mi := &file_proto_kademlia_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Nodes.ProtoReflect.Descriptor instead.
func (*Nodes) Descriptor() ([]byte, []int) {
	return file_proto_kademlia_proto_rawDescGZIP(), []int{4}
}

func (x *Nodes) GetNode() []*Node {
	if x != nil {
		return x.Node
	}
	return nil
}

type NodesOrData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Value:
	//
	//	*NodesOrData_Nodes
	//	*NodesOrData_Data
	Value isNodesOrData_Value `protobuf_oneof:"Value"`
}

func (x *NodesOrData) Reset() {
	*x = NodesOrData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_kademlia_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodesOrData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodesOrData) ProtoMessage() {}

func (x *NodesOrData) ProtoReflect() protoreflect.Message {
	mi := &file_proto_kademlia_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodesOrData.ProtoReflect.Descriptor instead.
func (*NodesOrData) Descriptor() ([]byte, []int) {
	return file_proto_kademlia_proto_rawDescGZIP(), []int{5}
}

func (m *NodesOrData) GetValue() isNodesOrData_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *NodesOrData) GetNodes() *Nodes {
	if x, ok := x.GetValue().(*NodesOrData_Nodes); ok {
		return x.Nodes
	}
	return nil
}

func (x *NodesOrData) GetData() string {
	if x, ok := x.GetValue().(*NodesOrData_Data); ok {
		return x.Data
	}
	return ""
}

type isNodesOrData_Value interface {
	isNodesOrData_Value()
}

type NodesOrData_Nodes struct {
	Nodes *Nodes `protobuf:"bytes,1,opt,name=Nodes,proto3,oneof"`
}

type NodesOrData_Data struct {
	Data string `protobuf:"bytes,2,opt,name=Data,proto3,oneof"`
}

func (*NodesOrData_Nodes) isNodesOrData_Value() {}

func (*NodesOrData_Data) isNodesOrData_Value() {}

var File_proto_kademlia_proto protoreflect.FileDescriptor

var file_proto_kademlia_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6b, 0x61, 0x64, 0x65, 0x6d, 0x6c, 0x69, 0x61,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x27, 0x0a,
	0x0b, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x18, 0x0a, 0x07,
	0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x53,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x44, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x12, 0x23, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x61, 0x64, 0x65, 0x6d, 0x6c, 0x69, 0x61, 0x49,
	0x44, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x22, 0x0a, 0x0a,
	0x4b, 0x61, 0x64, 0x65, 0x6d, 0x6c, 0x69, 0x61, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x49, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x21, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x61, 0x64,
	0x65, 0x6d, 0x6c, 0x69, 0x61, 0x49, 0x44, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x49,
	0x50, 0x57, 0x69, 0x74, 0x68, 0x50, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x49, 0x50, 0x57, 0x69, 0x74, 0x68, 0x50, 0x6f, 0x72, 0x74, 0x22, 0x28, 0x0a, 0x05, 0x4e,
	0x6f, 0x64, 0x65, 0x73, 0x12, 0x1f, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52,
	0x04, 0x4e, 0x6f, 0x64, 0x65, 0x22, 0x52, 0x0a, 0x0b, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x4f, 0x72,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x24, 0x0a, 0x05, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x64, 0x65,
	0x73, 0x48, 0x00, 0x52, 0x05, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x14, 0x0a, 0x04, 0x44, 0x61,
	0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61,
	0x42, 0x07, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x32, 0xc4, 0x01, 0x0a, 0x08, 0x4b, 0x61,
	0x64, 0x65, 0x6d, 0x6c, 0x69, 0x61, 0x12, 0x22, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x0b,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x1a, 0x0b, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x22, 0x00, 0x12, 0x2e, 0x0a, 0x09, 0x46, 0x69,
	0x6e, 0x64, 0x5f, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x4b, 0x61, 0x64, 0x65, 0x6d, 0x6c, 0x69, 0x61, 0x49, 0x44, 0x1a, 0x0c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x22, 0x00, 0x12, 0x35, 0x0a, 0x0a, 0x46, 0x69,
	0x6e, 0x64, 0x5f, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4b, 0x61, 0x64, 0x65, 0x6d, 0x6c, 0x69, 0x61, 0x49, 0x44, 0x1a, 0x12, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x4f, 0x72, 0x44, 0x61, 0x74, 0x61, 0x22,
	0x00, 0x12, 0x2d, 0x0a, 0x05, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x0e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x1a, 0x12, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00,
	0x42, 0x16, 0x5a, 0x14, 0x64, 0x37, 0x30, 0x32, 0x34, 0x65, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70,
	0x30, 0x34, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_kademlia_proto_rawDescOnce sync.Once
	file_proto_kademlia_proto_rawDescData = file_proto_kademlia_proto_rawDesc
)

func file_proto_kademlia_proto_rawDescGZIP() []byte {
	file_proto_kademlia_proto_rawDescOnce.Do(func() {
		file_proto_kademlia_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_kademlia_proto_rawDescData)
	})
	return file_proto_kademlia_proto_rawDescData
}

var file_proto_kademlia_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_kademlia_proto_goTypes = []any{
	(*StoreResult)(nil), // 0: proto.StoreResult
	(*Content)(nil),     // 1: proto.Content
	(*KademliaID)(nil),  // 2: proto.KademliaID
	(*Node)(nil),        // 3: proto.Node
	(*Nodes)(nil),       // 4: proto.Nodes
	(*NodesOrData)(nil), // 5: proto.NodesOrData
}
var file_proto_kademlia_proto_depIdxs = []int32{
	2, // 0: proto.Content.Key:type_name -> proto.KademliaID
	2, // 1: proto.Node.ID:type_name -> proto.KademliaID
	3, // 2: proto.Nodes.Node:type_name -> proto.Node
	4, // 3: proto.NodesOrData.Nodes:type_name -> proto.Nodes
	3, // 4: proto.Kademlia.Ping:input_type -> proto.Node
	2, // 5: proto.Kademlia.Find_Node:input_type -> proto.KademliaID
	2, // 6: proto.Kademlia.Find_Value:input_type -> proto.KademliaID
	1, // 7: proto.Kademlia.Store:input_type -> proto.Content
	3, // 8: proto.Kademlia.Ping:output_type -> proto.Node
	4, // 9: proto.Kademlia.Find_Node:output_type -> proto.Nodes
	5, // 10: proto.Kademlia.Find_Value:output_type -> proto.NodesOrData
	0, // 11: proto.Kademlia.Store:output_type -> proto.StoreResult
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_kademlia_proto_init() }
func file_proto_kademlia_proto_init() {
	if File_proto_kademlia_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_kademlia_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*StoreResult); i {
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
		file_proto_kademlia_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Content); i {
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
		file_proto_kademlia_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*KademliaID); i {
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
		file_proto_kademlia_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Node); i {
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
		file_proto_kademlia_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*Nodes); i {
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
		file_proto_kademlia_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*NodesOrData); i {
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
	file_proto_kademlia_proto_msgTypes[5].OneofWrappers = []any{
		(*NodesOrData_Nodes)(nil),
		(*NodesOrData_Data)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_kademlia_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_kademlia_proto_goTypes,
		DependencyIndexes: file_proto_kademlia_proto_depIdxs,
		MessageInfos:      file_proto_kademlia_proto_msgTypes,
	}.Build()
	File_proto_kademlia_proto = out.File
	file_proto_kademlia_proto_rawDesc = nil
	file_proto_kademlia_proto_goTypes = nil
	file_proto_kademlia_proto_depIdxs = nil
}
