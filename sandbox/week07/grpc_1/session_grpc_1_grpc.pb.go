// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: session_grpc_1.proto

// https://github.com/protocolbuffers/protobuf/releases/tag/v26.1
// https://github.com/protocolbuffers/protobuf/releases/download/v26.1/protoc-26.1-linux-x86_64.zip
// c:\bin\protoc-26.1-linux-x86_64\bin\protoc
// export PATH=${PATH}:/mnt/c/bin/protoc-26.1-linux-x86_64/bin:${HOME}/go/bin
// pushd sandbox/week07/grpc_1/
// go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
// go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
// go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
// protoc --go-grpc_out=. *.proto # protoc --go_out=plugins=grpc:. *.proto
// protoc --go_out=. --go-grpc_out=. *.proto

package grpc_1

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
	AuthChecker_Create_FullMethodName = "/grpc_1.AuthChecker/Create"
	AuthChecker_Check_FullMethodName  = "/grpc_1.AuthChecker/Check"
	AuthChecker_Delete_FullMethodName = "/grpc_1.AuthChecker/Delete"
)

// AuthCheckerClient is the client API for AuthChecker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthCheckerClient interface {
	Create(ctx context.Context, in *Session, opts ...grpc.CallOption) (*SessionID, error)
	Check(ctx context.Context, in *SessionID, opts ...grpc.CallOption) (*Session, error)
	Delete(ctx context.Context, in *SessionID, opts ...grpc.CallOption) (*Nothing, error)
}

type authCheckerClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthCheckerClient(cc grpc.ClientConnInterface) AuthCheckerClient {
	return &authCheckerClient{cc}
}

func (c *authCheckerClient) Create(ctx context.Context, in *Session, opts ...grpc.CallOption) (*SessionID, error) {
	out := new(SessionID)
	err := c.cc.Invoke(ctx, AuthChecker_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authCheckerClient) Check(ctx context.Context, in *SessionID, opts ...grpc.CallOption) (*Session, error) {
	out := new(Session)
	err := c.cc.Invoke(ctx, AuthChecker_Check_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authCheckerClient) Delete(ctx context.Context, in *SessionID, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := c.cc.Invoke(ctx, AuthChecker_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthCheckerServer is the server API for AuthChecker service.
// All implementations must embed UnimplementedAuthCheckerServer
// for forward compatibility
type AuthCheckerServer interface {
	Create(context.Context, *Session) (*SessionID, error)
	Check(context.Context, *SessionID) (*Session, error)
	Delete(context.Context, *SessionID) (*Nothing, error)
	mustEmbedUnimplementedAuthCheckerServer()
}

// UnimplementedAuthCheckerServer must be embedded to have forward compatible implementations.
type UnimplementedAuthCheckerServer struct {
}

func (UnimplementedAuthCheckerServer) Create(context.Context, *Session) (*SessionID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedAuthCheckerServer) Check(context.Context, *SessionID) (*Session, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedAuthCheckerServer) Delete(context.Context, *SessionID) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedAuthCheckerServer) mustEmbedUnimplementedAuthCheckerServer() {}

// UnsafeAuthCheckerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthCheckerServer will
// result in compilation errors.
type UnsafeAuthCheckerServer interface {
	mustEmbedUnimplementedAuthCheckerServer()
}

func RegisterAuthCheckerServer(s grpc.ServiceRegistrar, srv AuthCheckerServer) {
	s.RegisterService(&AuthChecker_ServiceDesc, srv)
}

func _AuthChecker_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Session)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthCheckerServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthChecker_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthCheckerServer).Create(ctx, req.(*Session))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthChecker_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthCheckerServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthChecker_Check_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthCheckerServer).Check(ctx, req.(*SessionID))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthChecker_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthCheckerServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthChecker_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthCheckerServer).Delete(ctx, req.(*SessionID))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthChecker_ServiceDesc is the grpc.ServiceDesc for AuthChecker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthChecker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc_1.AuthChecker",
	HandlerType: (*AuthCheckerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _AuthChecker_Create_Handler,
		},
		{
			MethodName: "Check",
			Handler:    _AuthChecker_Check_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _AuthChecker_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "session_grpc_1.proto",
}
