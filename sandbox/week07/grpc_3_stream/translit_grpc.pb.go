// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: translit.proto

// export PATH=${PATH}:/mnt/c/bin/protoc-26.1-linux-x86_64/bin:${HOME}/go/bin
// pushd sandbox/week07/grpc_3_stream
// protoc --go_out=. --go-grpc_out=. *.proto

package grpc_3_stream

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
	Transliteration_EnRu_FullMethodName = "/grpc_3_stream.Transliteration/EnRu"
)

// TransliterationClient is the client API for Transliteration service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TransliterationClient interface {
	EnRu(ctx context.Context, opts ...grpc.CallOption) (Transliteration_EnRuClient, error)
}

type transliterationClient struct {
	cc grpc.ClientConnInterface
}

func NewTransliterationClient(cc grpc.ClientConnInterface) TransliterationClient {
	return &transliterationClient{cc}
}

func (c *transliterationClient) EnRu(ctx context.Context, opts ...grpc.CallOption) (Transliteration_EnRuClient, error) {
	stream, err := c.cc.NewStream(ctx, &Transliteration_ServiceDesc.Streams[0], Transliteration_EnRu_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &transliterationEnRuClient{stream}
	return x, nil
}

type Transliteration_EnRuClient interface {
	Send(*Word) error
	Recv() (*Word, error)
	grpc.ClientStream
}

type transliterationEnRuClient struct {
	grpc.ClientStream
}

func (x *transliterationEnRuClient) Send(m *Word) error {
	return x.ClientStream.SendMsg(m)
}

func (x *transliterationEnRuClient) Recv() (*Word, error) {
	m := new(Word)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TransliterationServer is the server API for Transliteration service.
// All implementations must embed UnimplementedTransliterationServer
// for forward compatibility
type TransliterationServer interface {
	EnRu(Transliteration_EnRuServer) error
	mustEmbedUnimplementedTransliterationServer()
}

// UnimplementedTransliterationServer must be embedded to have forward compatible implementations.
type UnimplementedTransliterationServer struct {
}

func (UnimplementedTransliterationServer) EnRu(Transliteration_EnRuServer) error {
	return status.Errorf(codes.Unimplemented, "method EnRu not implemented")
}
func (UnimplementedTransliterationServer) mustEmbedUnimplementedTransliterationServer() {}

// UnsafeTransliterationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TransliterationServer will
// result in compilation errors.
type UnsafeTransliterationServer interface {
	mustEmbedUnimplementedTransliterationServer()
}

func RegisterTransliterationServer(s grpc.ServiceRegistrar, srv TransliterationServer) {
	s.RegisterService(&Transliteration_ServiceDesc, srv)
}

func _Transliteration_EnRu_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TransliterationServer).EnRu(&transliterationEnRuServer{stream})
}

type Transliteration_EnRuServer interface {
	Send(*Word) error
	Recv() (*Word, error)
	grpc.ServerStream
}

type transliterationEnRuServer struct {
	grpc.ServerStream
}

func (x *transliterationEnRuServer) Send(m *Word) error {
	return x.ServerStream.SendMsg(m)
}

func (x *transliterationEnRuServer) Recv() (*Word, error) {
	m := new(Word)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Transliteration_ServiceDesc is the grpc.ServiceDesc for Transliteration service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Transliteration_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc_3_stream.Transliteration",
	HandlerType: (*TransliterationServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "EnRu",
			Handler:       _Transliteration_EnRu_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "translit.proto",
}
