// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: proto/translator_user.proto

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

// TranslatorUserClient is the client API for TranslatorUser service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TranslatorUserClient interface {
	DictionaryLookup(ctx context.Context, in *DictionaryLookupParameter, opts ...grpc.CallOption) (*DictionaryLookupResposne, error)
	DictionaryLookupWithPos(ctx context.Context, in *DictionaryLookupWithPosParameter, opts ...grpc.CallOption) (*DictionaryLookupResposne, error)
}

type translatorUserClient struct {
	cc grpc.ClientConnInterface
}

func NewTranslatorUserClient(cc grpc.ClientConnInterface) TranslatorUserClient {
	return &translatorUserClient{cc}
}

func (c *translatorUserClient) DictionaryLookup(ctx context.Context, in *DictionaryLookupParameter, opts ...grpc.CallOption) (*DictionaryLookupResposne, error) {
	out := new(DictionaryLookupResposne)
	err := c.cc.Invoke(ctx, "/proto.TranslatorUser/DictionaryLookup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *translatorUserClient) DictionaryLookupWithPos(ctx context.Context, in *DictionaryLookupWithPosParameter, opts ...grpc.CallOption) (*DictionaryLookupResposne, error) {
	out := new(DictionaryLookupResposne)
	err := c.cc.Invoke(ctx, "/proto.TranslatorUser/DictionaryLookupWithPos", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TranslatorUserServer is the server API for TranslatorUser service.
// All implementations must embed UnimplementedTranslatorUserServer
// for forward compatibility
type TranslatorUserServer interface {
	DictionaryLookup(context.Context, *DictionaryLookupParameter) (*DictionaryLookupResposne, error)
	DictionaryLookupWithPos(context.Context, *DictionaryLookupWithPosParameter) (*DictionaryLookupResposne, error)
	mustEmbedUnimplementedTranslatorUserServer()
}

// UnimplementedTranslatorUserServer must be embedded to have forward compatible implementations.
type UnimplementedTranslatorUserServer struct {
}

func (UnimplementedTranslatorUserServer) DictionaryLookup(context.Context, *DictionaryLookupParameter) (*DictionaryLookupResposne, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DictionaryLookup not implemented")
}
func (UnimplementedTranslatorUserServer) DictionaryLookupWithPos(context.Context, *DictionaryLookupWithPosParameter) (*DictionaryLookupResposne, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DictionaryLookupWithPos not implemented")
}
func (UnimplementedTranslatorUserServer) mustEmbedUnimplementedTranslatorUserServer() {}

// UnsafeTranslatorUserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TranslatorUserServer will
// result in compilation errors.
type UnsafeTranslatorUserServer interface {
	mustEmbedUnimplementedTranslatorUserServer()
}

func RegisterTranslatorUserServer(s grpc.ServiceRegistrar, srv TranslatorUserServer) {
	s.RegisterService(&TranslatorUser_ServiceDesc, srv)
}

func _TranslatorUser_DictionaryLookup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DictionaryLookupParameter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TranslatorUserServer).DictionaryLookup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TranslatorUser/DictionaryLookup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TranslatorUserServer).DictionaryLookup(ctx, req.(*DictionaryLookupParameter))
	}
	return interceptor(ctx, in, info, handler)
}

func _TranslatorUser_DictionaryLookupWithPos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DictionaryLookupWithPosParameter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TranslatorUserServer).DictionaryLookupWithPos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TranslatorUser/DictionaryLookupWithPos",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TranslatorUserServer).DictionaryLookupWithPos(ctx, req.(*DictionaryLookupWithPosParameter))
	}
	return interceptor(ctx, in, info, handler)
}

// TranslatorUser_ServiceDesc is the grpc.ServiceDesc for TranslatorUser service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TranslatorUser_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.TranslatorUser",
	HandlerType: (*TranslatorUserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DictionaryLookup",
			Handler:    _TranslatorUser_DictionaryLookup_Handler,
		},
		{
			MethodName: "DictionaryLookupWithPos",
			Handler:    _TranslatorUser_DictionaryLookupWithPos_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/translator_user.proto",
}
