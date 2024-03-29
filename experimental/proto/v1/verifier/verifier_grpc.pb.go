// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: verifier.proto

package verifier

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

// VerifierPluginClient is the client API for VerifierPlugin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VerifierPluginClient interface {
	// Perform verification of a given artifact
	VerifyReference(ctx context.Context, in *VerifyReferenceRequest, opts ...grpc.CallOption) (*VerifyReferenceResponse, error)
}

type verifierPluginClient struct {
	cc grpc.ClientConnInterface
}

func NewVerifierPluginClient(cc grpc.ClientConnInterface) VerifierPluginClient {
	return &verifierPluginClient{cc}
}

func (c *verifierPluginClient) VerifyReference(ctx context.Context, in *VerifyReferenceRequest, opts ...grpc.CallOption) (*VerifyReferenceResponse, error) {
	out := new(VerifyReferenceResponse)
	err := c.cc.Invoke(ctx, "/verifier.VerifierPlugin/VerifyReference", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VerifierPluginServer is the server API for VerifierPlugin service.
// All implementations must embed UnimplementedVerifierPluginServer
// for forward compatibility
type VerifierPluginServer interface {
	// Perform verification of a given artifact
	VerifyReference(context.Context, *VerifyReferenceRequest) (*VerifyReferenceResponse, error)
	mustEmbedUnimplementedVerifierPluginServer()
}

// UnimplementedVerifierPluginServer must be embedded to have forward compatible implementations.
type UnimplementedVerifierPluginServer struct {
}

func (UnimplementedVerifierPluginServer) VerifyReference(context.Context, *VerifyReferenceRequest) (*VerifyReferenceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyReference not implemented")
}
func (UnimplementedVerifierPluginServer) mustEmbedUnimplementedVerifierPluginServer() {}

// UnsafeVerifierPluginServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VerifierPluginServer will
// result in compilation errors.
type UnsafeVerifierPluginServer interface {
	mustEmbedUnimplementedVerifierPluginServer()
}

func RegisterVerifierPluginServer(s grpc.ServiceRegistrar, srv VerifierPluginServer) {
	s.RegisterService(&VerifierPlugin_ServiceDesc, srv)
}

func _VerifierPlugin_VerifyReference_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyReferenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VerifierPluginServer).VerifyReference(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/verifier.VerifierPlugin/VerifyReference",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VerifierPluginServer).VerifyReference(ctx, req.(*VerifyReferenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VerifierPlugin_ServiceDesc is the grpc.ServiceDesc for VerifierPlugin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VerifierPlugin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "verifier.VerifierPlugin",
	HandlerType: (*VerifierPluginServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "VerifyReference",
			Handler:    _VerifierPlugin_VerifyReference_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "verifier.proto",
}
