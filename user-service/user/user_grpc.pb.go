// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package u

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserClient interface {
	FindAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (User_FindAllClient, error)
	SaveUser(ctx context.Context, in *SaveUserRequest, opts ...grpc.CallOption) (*SaveUserResponse, error)
}

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

func (c *userClient) FindAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (User_FindAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &User_ServiceDesc.Streams[0], "/user.User/FindAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &userFindAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type User_FindAllClient interface {
	Recv() (*UserResponse, error)
	grpc.ClientStream
}

type userFindAllClient struct {
	grpc.ClientStream
}

func (x *userFindAllClient) Recv() (*UserResponse, error) {
	m := new(UserResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *userClient) SaveUser(ctx context.Context, in *SaveUserRequest, opts ...grpc.CallOption) (*SaveUserResponse, error) {
	out := new(SaveUserResponse)
	err := c.cc.Invoke(ctx, "/user.User/SaveUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServer is the server API for User service.
// All implementations must embed UnimplementedUserServer
// for forward compatibility
type UserServer interface {
	FindAll(*emptypb.Empty, User_FindAllServer) error
	SaveUser(context.Context, *SaveUserRequest) (*SaveUserResponse, error)
	mustEmbedUnimplementedUserServer()
}

// UnimplementedUserServer must be embedded to have forward compatible implementations.
type UnimplementedUserServer struct {
}

func (UnimplementedUserServer) FindAll(*emptypb.Empty, User_FindAllServer) error {
	return status.Errorf(codes.Unimplemented, "method FindAll not implemented")
}
func (UnimplementedUserServer) SaveUser(context.Context, *SaveUserRequest) (*SaveUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveUser not implemented")
}
func (UnimplementedUserServer) mustEmbedUnimplementedUserServer() {}

// UnsafeUserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServer will
// result in compilation errors.
type UnsafeUserServer interface {
	mustEmbedUnimplementedUserServer()
}

func RegisterUserServer(s grpc.ServiceRegistrar, srv UserServer) {
	s.RegisterService(&User_ServiceDesc, srv)
}

func _User_FindAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UserServer).FindAll(m, &userFindAllServer{stream})
}

type User_FindAllServer interface {
	Send(*UserResponse) error
	grpc.ServerStream
}

type userFindAllServer struct {
	grpc.ServerStream
}

func (x *userFindAllServer) Send(m *UserResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _User_SaveUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).SaveUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/SaveUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).SaveUser(ctx, req.(*SaveUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// User_ServiceDesc is the grpc.ServiceDesc for User service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var User_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.User",
	HandlerType: (*UserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SaveUser",
			Handler:    _User_SaveUser_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "FindAll",
			Handler:       _User_FindAll_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "user.proto",
}
