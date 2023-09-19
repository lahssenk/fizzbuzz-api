// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: fizzbuzz/v1/service.proto

package fizzbuzz_v1

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

const (
	FizzBuzzService_ComputeFizzBuzzRange_FullMethodName = "/fizzbuzz.v1.FizzBuzzService/ComputeFizzBuzzRange"
)

// FizzBuzzServiceClient is the client API for FizzBuzzService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FizzBuzzServiceClient interface {
	// Compute the FizzBuzz value for a range of integers, from 1 to <limit>
	ComputeFizzBuzzRange(ctx context.Context, in *ComputeFizzBuzzRangeRequest, opts ...grpc.CallOption) (*ComputeFizzBuzzRangeResponse, error)
}

type fizzBuzzServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFizzBuzzServiceClient(cc grpc.ClientConnInterface) FizzBuzzServiceClient {
	return &fizzBuzzServiceClient{cc}
}

func (c *fizzBuzzServiceClient) ComputeFizzBuzzRange(ctx context.Context, in *ComputeFizzBuzzRangeRequest, opts ...grpc.CallOption) (*ComputeFizzBuzzRangeResponse, error) {
	out := new(ComputeFizzBuzzRangeResponse)
	err := c.cc.Invoke(ctx, FizzBuzzService_ComputeFizzBuzzRange_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FizzBuzzServiceServer is the server API for FizzBuzzService service.
// All implementations must embed UnimplementedFizzBuzzServiceServer
// for forward compatibility
type FizzBuzzServiceServer interface {
	// Compute the FizzBuzz value for a range of integers, from 1 to <limit>
	ComputeFizzBuzzRange(context.Context, *ComputeFizzBuzzRangeRequest) (*ComputeFizzBuzzRangeResponse, error)
	mustEmbedUnimplementedFizzBuzzServiceServer()
}

// UnimplementedFizzBuzzServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFizzBuzzServiceServer struct {
}

func (UnimplementedFizzBuzzServiceServer) ComputeFizzBuzzRange(context.Context, *ComputeFizzBuzzRangeRequest) (*ComputeFizzBuzzRangeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ComputeFizzBuzzRange not implemented")
}
func (UnimplementedFizzBuzzServiceServer) mustEmbedUnimplementedFizzBuzzServiceServer() {}

// UnsafeFizzBuzzServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FizzBuzzServiceServer will
// result in compilation errors.
type UnsafeFizzBuzzServiceServer interface {
	mustEmbedUnimplementedFizzBuzzServiceServer()
}

func RegisterFizzBuzzServiceServer(s grpc.ServiceRegistrar, srv FizzBuzzServiceServer) {
	s.RegisterService(&FizzBuzzService_ServiceDesc, srv)
}

func _FizzBuzzService_ComputeFizzBuzzRange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ComputeFizzBuzzRangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FizzBuzzServiceServer).ComputeFizzBuzzRange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FizzBuzzService_ComputeFizzBuzzRange_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FizzBuzzServiceServer).ComputeFizzBuzzRange(ctx, req.(*ComputeFizzBuzzRangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FizzBuzzService_ServiceDesc is the grpc.ServiceDesc for FizzBuzzService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FizzBuzzService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fizzbuzz.v1.FizzBuzzService",
	HandlerType: (*FizzBuzzServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ComputeFizzBuzzRange",
			Handler:    _FizzBuzzService_ComputeFizzBuzzRange_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fizzbuzz/v1/service.proto",
}