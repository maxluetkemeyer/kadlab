// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: proto/kademlia.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Kademlia_Ping_FullMethodName       = "/proto.Kademlia/Ping"
	Kademlia_Find_Node_FullMethodName  = "/proto.Kademlia/Find_Node"
	Kademlia_Find_Value_FullMethodName = "/proto.Kademlia/Find_Value"
	Kademlia_Store_FullMethodName      = "/proto.Kademlia/Store"
)

// KademliaClient is the client API for Kademlia service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KademliaClient interface {
	Ping(ctx context.Context, in *Node, opts ...grpc.CallOption) (*Node, error)
	Find_Node(ctx context.Context, in *KademliaID, opts ...grpc.CallOption) (*Nodes, error)
	Find_Value(ctx context.Context, in *KademliaID, opts ...grpc.CallOption) (*NodesOrData, error)
	Store(ctx context.Context, in *Content, opts ...grpc.CallOption) (*StoreResult, error)
}

type kademliaClient struct {
	cc grpc.ClientConnInterface
}

func NewKademliaClient(cc grpc.ClientConnInterface) KademliaClient {
	return &kademliaClient{cc}
}

func (c *kademliaClient) Ping(ctx context.Context, in *Node, opts ...grpc.CallOption) (*Node, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Node)
	err := c.cc.Invoke(ctx, Kademlia_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kademliaClient) Find_Node(ctx context.Context, in *KademliaID, opts ...grpc.CallOption) (*Nodes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Nodes)
	err := c.cc.Invoke(ctx, Kademlia_Find_Node_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kademliaClient) Find_Value(ctx context.Context, in *KademliaID, opts ...grpc.CallOption) (*NodesOrData, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NodesOrData)
	err := c.cc.Invoke(ctx, Kademlia_Find_Value_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kademliaClient) Store(ctx context.Context, in *Content, opts ...grpc.CallOption) (*StoreResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StoreResult)
	err := c.cc.Invoke(ctx, Kademlia_Store_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KademliaServer is the server API for Kademlia service.
// All implementations must embed UnimplementedKademliaServer
// for forward compatibility.
type KademliaServer interface {
	Ping(context.Context, *Node) (*Node, error)
	Find_Node(context.Context, *KademliaID) (*Nodes, error)
	Find_Value(context.Context, *KademliaID) (*NodesOrData, error)
	Store(context.Context, *Content) (*StoreResult, error)
	mustEmbedUnimplementedKademliaServer()
}

// UnimplementedKademliaServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedKademliaServer struct{}

func (UnimplementedKademliaServer) Ping(context.Context, *Node) (*Node, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedKademliaServer) Find_Node(context.Context, *KademliaID) (*Nodes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find_Node not implemented")
}
func (UnimplementedKademliaServer) Find_Value(context.Context, *KademliaID) (*NodesOrData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find_Value not implemented")
}
func (UnimplementedKademliaServer) Store(context.Context, *Content) (*StoreResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Store not implemented")
}
func (UnimplementedKademliaServer) mustEmbedUnimplementedKademliaServer() {}
func (UnimplementedKademliaServer) testEmbeddedByValue()                  {}

// UnsafeKademliaServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KademliaServer will
// result in compilation errors.
type UnsafeKademliaServer interface {
	mustEmbedUnimplementedKademliaServer()
}

func RegisterKademliaServer(s grpc.ServiceRegistrar, srv KademliaServer) {
	// If the following call pancis, it indicates UnimplementedKademliaServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Kademlia_ServiceDesc, srv)
}

func _Kademlia_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Node)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KademliaServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Kademlia_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KademliaServer).Ping(ctx, req.(*Node))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kademlia_Find_Node_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KademliaID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KademliaServer).Find_Node(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Kademlia_Find_Node_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KademliaServer).Find_Node(ctx, req.(*KademliaID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kademlia_Find_Value_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KademliaID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KademliaServer).Find_Value(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Kademlia_Find_Value_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KademliaServer).Find_Value(ctx, req.(*KademliaID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kademlia_Store_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Content)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KademliaServer).Store(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Kademlia_Store_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KademliaServer).Store(ctx, req.(*Content))
	}
	return interceptor(ctx, in, info, handler)
}

// Kademlia_ServiceDesc is the grpc.ServiceDesc for Kademlia service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Kademlia_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Kademlia",
	HandlerType: (*KademliaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Kademlia_Ping_Handler,
		},
		{
			MethodName: "Find_Node",
			Handler:    _Kademlia_Find_Node_Handler,
		},
		{
			MethodName: "Find_Value",
			Handler:    _Kademlia_Find_Value_Handler,
		},
		{
			MethodName: "Store",
			Handler:    _Kademlia_Store_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/kademlia.proto",
}
