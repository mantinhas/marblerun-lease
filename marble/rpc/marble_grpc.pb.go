// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package rpc

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

// OperationTrackerClient is the client API for OperationTracker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OperationTrackerClient interface {
	StartOperation(ctx context.Context, in *StartOperationReq, opts ...grpc.CallOption) (*StartOperationResp, error)
}

type operationTrackerClient struct {
	cc grpc.ClientConnInterface
}

func NewOperationTrackerClient(cc grpc.ClientConnInterface) OperationTrackerClient {
	return &operationTrackerClient{cc}
}

func (c *operationTrackerClient) StartOperation(ctx context.Context, in *StartOperationReq, opts ...grpc.CallOption) (*StartOperationResp, error) {
	out := new(StartOperationResp)
	err := c.cc.Invoke(ctx, "/rpc.OperationTracker/StartOperation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OperationTrackerServer is the server API for OperationTracker service.
// All implementations must embed UnimplementedOperationTrackerServer
// for forward compatibility
type OperationTrackerServer interface {
	StartOperation(context.Context, *StartOperationReq) (*StartOperationResp, error)
	mustEmbedUnimplementedOperationTrackerServer()
}

// UnimplementedOperationTrackerServer must be embedded to have forward compatible implementations.
type UnimplementedOperationTrackerServer struct {
}

func (UnimplementedOperationTrackerServer) StartOperation(context.Context, *StartOperationReq) (*StartOperationResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartOperation not implemented")
}
func (UnimplementedOperationTrackerServer) mustEmbedUnimplementedOperationTrackerServer() {}

// UnsafeOperationTrackerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OperationTrackerServer will
// result in compilation errors.
type UnsafeOperationTrackerServer interface {
	mustEmbedUnimplementedOperationTrackerServer()
}

func RegisterOperationTrackerServer(s grpc.ServiceRegistrar, srv OperationTrackerServer) {
	s.RegisterService(&OperationTracker_ServiceDesc, srv)
}

func _OperationTracker_StartOperation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartOperationReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationTrackerServer).StartOperation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.OperationTracker/StartOperation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationTrackerServer).StartOperation(ctx, req.(*StartOperationReq))
	}
	return interceptor(ctx, in, info, handler)
}

// OperationTracker_ServiceDesc is the grpc.ServiceDesc for OperationTracker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OperationTracker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.OperationTracker",
	HandlerType: (*OperationTrackerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartOperation",
			Handler:    _OperationTracker_StartOperation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "marble.proto",
}
