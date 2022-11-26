// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package sdk

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BabelServiceClient is the client API for BabelService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BabelServiceClient interface {
	GetSyncStatus(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*SyncStatus, error)
}

type babelServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBabelServiceClient(cc grpc.ClientConnInterface) BabelServiceClient {
	return &babelServiceClient{cc}
}

func (c *babelServiceClient) GetSyncStatus(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*SyncStatus, error) {
	out := new(SyncStatus)
	err := c.cc.Invoke(ctx, "/proto.BabelService/GetSyncStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BabelServiceServer is the server API for BabelService service.
// All implementations must embed UnimplementedBabelServiceServer
// for forward compatibility
type BabelServiceServer interface {
	GetSyncStatus(context.Context, *empty.Empty) (*SyncStatus, error)
	mustEmbedUnimplementedBabelServiceServer()
}

// UnimplementedBabelServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBabelServiceServer struct {
}

func (UnimplementedBabelServiceServer) GetSyncStatus(context.Context, *empty.Empty) (*SyncStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSyncStatus not implemented")
}
func (UnimplementedBabelServiceServer) mustEmbedUnimplementedBabelServiceServer() {}

// UnsafeBabelServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BabelServiceServer will
// result in compilation errors.
type UnsafeBabelServiceServer interface {
	mustEmbedUnimplementedBabelServiceServer()
}

func RegisterBabelServiceServer(s grpc.ServiceRegistrar, srv BabelServiceServer) {
	s.RegisterService(&BabelService_ServiceDesc, srv)
}

func _BabelService_GetSyncStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BabelServiceServer).GetSyncStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.BabelService/GetSyncStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BabelServiceServer).GetSyncStatus(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// BabelService_ServiceDesc is the grpc.ServiceDesc for BabelService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BabelService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.BabelService",
	HandlerType: (*BabelServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSyncStatus",
			Handler:    _BabelService_GetSyncStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sdk/babel.proto",
}