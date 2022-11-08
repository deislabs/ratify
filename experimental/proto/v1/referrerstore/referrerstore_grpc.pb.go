// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: referrerstore.proto

package referrerstore

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ReferrerStorePluginClient is the client API for ReferrerStorePlugin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReferrerStorePluginClient interface {
	// Query all the referrers that are linked to a subject.
	ListReferrers(ctx context.Context, in *ListReferrersRequest, opts ...grpc.CallOption) (ReferrerStorePlugin_ListReferrersClient, error)
	// Fetch the contents of a blob for a given artifact.
	GetBlobContent(ctx context.Context, in *GetBlobContentRequest, opts ...grpc.CallOption) (*GetBlobContentResponse, error)
	// Fetch additional metadata for a subject.
	GetSubjectDescriptor(ctx context.Context, in *GetSubjectDescriptorRequest, opts ...grpc.CallOption) (*GetSubjectDescriptorResponse, error)
	// Fetch the contents of a reference manifest.
	GetReferenceManifest(ctx context.Context, in *GetManifestRequest, opts ...grpc.CallOption) (*GetManifestResponse, error)
}

type referrerStorePluginClient struct {
	cc grpc.ClientConnInterface
}

func NewReferrerStorePluginClient(cc grpc.ClientConnInterface) ReferrerStorePluginClient {
	return &referrerStorePluginClient{cc}
}

func (c *referrerStorePluginClient) ListReferrers(ctx context.Context, in *ListReferrersRequest, opts ...grpc.CallOption) (ReferrerStorePlugin_ListReferrersClient, error) {
	stream, err := c.cc.NewStream(ctx, &ReferrerStorePlugin_ServiceDesc.Streams[0], "/referrerstore.ReferrerStorePlugin/ListReferrers", opts...)
	if err != nil {
		return nil, err
	}
	x := &referrerStorePluginListReferrersClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ReferrerStorePlugin_ListReferrersClient interface {
	Recv() (*ListReferrersResponse, error)
	grpc.ClientStream
}

type referrerStorePluginListReferrersClient struct {
	grpc.ClientStream
}

func (x *referrerStorePluginListReferrersClient) Recv() (*ListReferrersResponse, error) {
	m := new(ListReferrersResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *referrerStorePluginClient) GetBlobContent(ctx context.Context, in *GetBlobContentRequest, opts ...grpc.CallOption) (*GetBlobContentResponse, error) {
	out := new(GetBlobContentResponse)
	err := c.cc.Invoke(ctx, "/referrerstore.ReferrerStorePlugin/GetBlobContent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *referrerStorePluginClient) GetSubjectDescriptor(ctx context.Context, in *GetSubjectDescriptorRequest, opts ...grpc.CallOption) (*GetSubjectDescriptorResponse, error) {
	out := new(GetSubjectDescriptorResponse)
	err := c.cc.Invoke(ctx, "/referrerstore.ReferrerStorePlugin/GetSubjectDescriptor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *referrerStorePluginClient) GetReferenceManifest(ctx context.Context, in *GetManifestRequest, opts ...grpc.CallOption) (*GetManifestResponse, error) {
	out := new(GetManifestResponse)
	err := c.cc.Invoke(ctx, "/referrerstore.ReferrerStorePlugin/GetReferenceManifest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReferrerStorePluginServer is the server API for ReferrerStorePlugin service.
// All implementations must embed UnimplementedReferrerStorePluginServer
// for forward compatibility
type ReferrerStorePluginServer interface {
	// Query all the referrers that are linked to a subject.
	ListReferrers(*ListReferrersRequest, ReferrerStorePlugin_ListReferrersServer) error
	// Fetch the contents of a blob for a given artifact.
	GetBlobContent(context.Context, *GetBlobContentRequest) (*GetBlobContentResponse, error)
	// Fetch additional metadata for a subject.
	GetSubjectDescriptor(context.Context, *GetSubjectDescriptorRequest) (*GetSubjectDescriptorResponse, error)
	// Fetch the contents of a reference manifest.
	GetReferenceManifest(context.Context, *GetManifestRequest) (*GetManifestResponse, error)
	mustEmbedUnimplementedReferrerStorePluginServer()
}

// UnimplementedReferrerStorePluginServer must be embedded to have forward compatible implementations.
type UnimplementedReferrerStorePluginServer struct {
}

func (UnimplementedReferrerStorePluginServer) ListReferrers(*ListReferrersRequest, ReferrerStorePlugin_ListReferrersServer) error {
	return status.Errorf(codes.Unimplemented, "method ListReferrers not implemented")
}
func (UnimplementedReferrerStorePluginServer) GetBlobContent(context.Context, *GetBlobContentRequest) (*GetBlobContentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlobContent not implemented")
}
func (UnimplementedReferrerStorePluginServer) GetSubjectDescriptor(context.Context, *GetSubjectDescriptorRequest) (*GetSubjectDescriptorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSubjectDescriptor not implemented")
}
func (UnimplementedReferrerStorePluginServer) GetReferenceManifest(context.Context, *GetManifestRequest) (*GetManifestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReferenceManifest not implemented")
}
func (UnimplementedReferrerStorePluginServer) mustEmbedUnimplementedReferrerStorePluginServer() {}

// UnsafeReferrerStorePluginServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReferrerStorePluginServer will
// result in compilation errors.
type UnsafeReferrerStorePluginServer interface {
	mustEmbedUnimplementedReferrerStorePluginServer()
}

func RegisterReferrerStorePluginServer(s grpc.ServiceRegistrar, srv ReferrerStorePluginServer) {
	s.RegisterService(&ReferrerStorePlugin_ServiceDesc, srv)
}

func _ReferrerStorePlugin_ListReferrers_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListReferrersRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ReferrerStorePluginServer).ListReferrers(m, &referrerStorePluginListReferrersServer{stream})
}

type ReferrerStorePlugin_ListReferrersServer interface {
	Send(*ListReferrersResponse) error
	grpc.ServerStream
}

type referrerStorePluginListReferrersServer struct {
	grpc.ServerStream
}

func (x *referrerStorePluginListReferrersServer) Send(m *ListReferrersResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _ReferrerStorePlugin_GetBlobContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlobContentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReferrerStorePluginServer).GetBlobContent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/referrerstore.ReferrerStorePlugin/GetBlobContent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReferrerStorePluginServer).GetBlobContent(ctx, req.(*GetBlobContentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReferrerStorePlugin_GetSubjectDescriptor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSubjectDescriptorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReferrerStorePluginServer).GetSubjectDescriptor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/referrerstore.ReferrerStorePlugin/GetSubjectDescriptor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReferrerStorePluginServer).GetSubjectDescriptor(ctx, req.(*GetSubjectDescriptorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReferrerStorePlugin_GetReferenceManifest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetManifestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReferrerStorePluginServer).GetReferenceManifest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/referrerstore.ReferrerStorePlugin/GetReferenceManifest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReferrerStorePluginServer).GetReferenceManifest(ctx, req.(*GetManifestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ReferrerStorePlugin_ServiceDesc is the grpc.ServiceDesc for ReferrerStorePlugin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReferrerStorePlugin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "referrerstore.ReferrerStorePlugin",
	HandlerType: (*ReferrerStorePluginServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBlobContent",
			Handler:    _ReferrerStorePlugin_GetBlobContent_Handler,
		},
		{
			MethodName: "GetSubjectDescriptor",
			Handler:    _ReferrerStorePlugin_GetSubjectDescriptor_Handler,
		},
		{
			MethodName: "GetReferenceManifest",
			Handler:    _ReferrerStorePlugin_GetReferenceManifest_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListReferrers",
			Handler:       _ReferrerStorePlugin_ListReferrers_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "referrerstore.proto",
}
