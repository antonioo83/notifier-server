// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.2
// source: grpc/message.proto

package message_proto

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

// MessageClient is the client API for Message service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageClient interface {
	CreateMessages(ctx context.Context, in *MessagesCreateRequest, opts ...grpc.CallOption) (*MessagesCreateResponse, error)
	DeleteMessage(ctx context.Context, in *MessageDeleteRequest, opts ...grpc.CallOption) (*MessageDeleteResponse, error)
	GetMessage(ctx context.Context, in *MessageGetRequest, opts ...grpc.CallOption) (*MessageGetResponse, error)
}

type messageClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageClient(cc grpc.ClientConnInterface) MessageClient {
	return &messageClient{cc}
}

func (c *messageClient) CreateMessages(ctx context.Context, in *MessagesCreateRequest, opts ...grpc.CallOption) (*MessagesCreateResponse, error) {
	out := new(MessagesCreateResponse)
	err := c.cc.Invoke(ctx, "/grpc.Message/CreateMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageClient) DeleteMessage(ctx context.Context, in *MessageDeleteRequest, opts ...grpc.CallOption) (*MessageDeleteResponse, error) {
	out := new(MessageDeleteResponse)
	err := c.cc.Invoke(ctx, "/grpc.Message/DeleteMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageClient) GetMessage(ctx context.Context, in *MessageGetRequest, opts ...grpc.CallOption) (*MessageGetResponse, error) {
	out := new(MessageGetResponse)
	err := c.cc.Invoke(ctx, "/grpc.Message/GetMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessageServer is the server API for Message service.
// All implementations must embed UnimplementedMessageServer
// for forward compatibility
type MessageServer interface {
	CreateMessages(context.Context, *MessagesCreateRequest) (*MessagesCreateResponse, error)
	DeleteMessage(context.Context, *MessageDeleteRequest) (*MessageDeleteResponse, error)
	GetMessage(context.Context, *MessageGetRequest) (*MessageGetResponse, error)
	mustEmbedUnimplementedMessageServer()
}

// UnimplementedMessageServer must be embedded to have forward compatible implementations.
type UnimplementedMessageServer struct {
}

func (UnimplementedMessageServer) CreateMessages(context.Context, *MessagesCreateRequest) (*MessagesCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMessages not implemented")
}
func (UnimplementedMessageServer) DeleteMessage(context.Context, *MessageDeleteRequest) (*MessageDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMessage not implemented")
}
func (UnimplementedMessageServer) GetMessage(context.Context, *MessageGetRequest) (*MessageGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMessage not implemented")
}
func (UnimplementedMessageServer) mustEmbedUnimplementedMessageServer() {}

// UnsafeMessageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServer will
// result in compilation errors.
type UnsafeMessageServer interface {
	mustEmbedUnimplementedMessageServer()
}

func RegisterMessageServer(s grpc.ServiceRegistrar, srv MessageServer) {
	s.RegisterService(&Message_ServiceDesc, srv)
}

func _Message_CreateMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessagesCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServer).CreateMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Message/CreateMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServer).CreateMessages(ctx, req.(*MessagesCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Message_DeleteMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServer).DeleteMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Message/DeleteMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServer).DeleteMessage(ctx, req.(*MessageDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Message_GetMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServer).GetMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Message/GetMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServer).GetMessage(ctx, req.(*MessageGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Message_ServiceDesc is the grpc.ServiceDesc for Message service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Message_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Message",
	HandlerType: (*MessageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMessages",
			Handler:    _Message_CreateMessages_Handler,
		},
		{
			MethodName: "DeleteMessage",
			Handler:    _Message_DeleteMessage_Handler,
		},
		{
			MethodName: "GetMessage",
			Handler:    _Message_GetMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/message.proto",
}
