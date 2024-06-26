// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.27.0
// source: service.proto

package proto

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
	PasswordVaultService_UserCreate_FullMethodName = "/grpc.PasswordVaultService/UserCreate"
	PasswordVaultService_UserLogin_FullMethodName  = "/grpc.PasswordVaultService/UserLogin"
	PasswordVaultService_DataWrite_FullMethodName  = "/grpc.PasswordVaultService/DataWrite"
	PasswordVaultService_DataRead_FullMethodName   = "/grpc.PasswordVaultService/DataRead"
)

// PasswordVaultServiceClient is the client API for PasswordVaultService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PasswordVaultServiceClient interface {
	UserCreate(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error)
	UserLogin(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error)
	DataWrite(ctx context.Context, in *DataWriteRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	DataRead(ctx context.Context, in *DataReadRequest, opts ...grpc.CallOption) (*DataReadResponse, error)
}

type passwordVaultServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPasswordVaultServiceClient(cc grpc.ClientConnInterface) PasswordVaultServiceClient {
	return &passwordVaultServiceClient{cc}
}

func (c *passwordVaultServiceClient) UserCreate(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, PasswordVaultService_UserCreate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *passwordVaultServiceClient) UserLogin(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, PasswordVaultService_UserLogin_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *passwordVaultServiceClient) DataWrite(ctx context.Context, in *DataWriteRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, PasswordVaultService_DataWrite_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *passwordVaultServiceClient) DataRead(ctx context.Context, in *DataReadRequest, opts ...grpc.CallOption) (*DataReadResponse, error) {
	out := new(DataReadResponse)
	err := c.cc.Invoke(ctx, PasswordVaultService_DataRead_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PasswordVaultServiceServer is the server API for PasswordVaultService service.
// All implementations must embed UnimplementedPasswordVaultServiceServer
// for forward compatibility
type PasswordVaultServiceServer interface {
	UserCreate(context.Context, *UserRequest) (*UserResponse, error)
	UserLogin(context.Context, *UserRequest) (*UserResponse, error)
	DataWrite(context.Context, *DataWriteRequest) (*EmptyResponse, error)
	DataRead(context.Context, *DataReadRequest) (*DataReadResponse, error)
	mustEmbedUnimplementedPasswordVaultServiceServer()
}

// UnimplementedPasswordVaultServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPasswordVaultServiceServer struct {
}

func (UnimplementedPasswordVaultServiceServer) UserCreate(context.Context, *UserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserCreate not implemented")
}
func (UnimplementedPasswordVaultServiceServer) UserLogin(context.Context, *UserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserLogin not implemented")
}
func (UnimplementedPasswordVaultServiceServer) DataWrite(context.Context, *DataWriteRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DataWrite not implemented")
}
func (UnimplementedPasswordVaultServiceServer) DataRead(context.Context, *DataReadRequest) (*DataReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DataRead not implemented")
}
func (UnimplementedPasswordVaultServiceServer) mustEmbedUnimplementedPasswordVaultServiceServer() {}

// UnsafePasswordVaultServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PasswordVaultServiceServer will
// result in compilation errors.
type UnsafePasswordVaultServiceServer interface {
	mustEmbedUnimplementedPasswordVaultServiceServer()
}

func RegisterPasswordVaultServiceServer(s grpc.ServiceRegistrar, srv PasswordVaultServiceServer) {
	s.RegisterService(&PasswordVaultService_ServiceDesc, srv)
}

func _PasswordVaultService_UserCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PasswordVaultServiceServer).UserCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PasswordVaultService_UserCreate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PasswordVaultServiceServer).UserCreate(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PasswordVaultService_UserLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PasswordVaultServiceServer).UserLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PasswordVaultService_UserLogin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PasswordVaultServiceServer).UserLogin(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PasswordVaultService_DataWrite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DataWriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PasswordVaultServiceServer).DataWrite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PasswordVaultService_DataWrite_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PasswordVaultServiceServer).DataWrite(ctx, req.(*DataWriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PasswordVaultService_DataRead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DataReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PasswordVaultServiceServer).DataRead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PasswordVaultService_DataRead_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PasswordVaultServiceServer).DataRead(ctx, req.(*DataReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PasswordVaultService_ServiceDesc is the grpc.ServiceDesc for PasswordVaultService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PasswordVaultService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.PasswordVaultService",
	HandlerType: (*PasswordVaultServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserCreate",
			Handler:    _PasswordVaultService_UserCreate_Handler,
		},
		{
			MethodName: "UserLogin",
			Handler:    _PasswordVaultService_UserLogin_Handler,
		},
		{
			MethodName: "DataWrite",
			Handler:    _PasswordVaultService_DataWrite_Handler,
		},
		{
			MethodName: "DataRead",
			Handler:    _PasswordVaultService_DataRead_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
