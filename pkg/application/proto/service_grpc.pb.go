// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: service.proto

package proto

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

const (
	ClerkAPIService_OpenAccount_FullMethodName   = "/ClerkAPIService/OpenAccount"
	ClerkAPIService_ListAccounts_FullMethodName  = "/ClerkAPIService/ListAccounts"
	ClerkAPIService_AddMoney_FullMethodName      = "/ClerkAPIService/AddMoney"
	ClerkAPIService_WithdrawMoney_FullMethodName = "/ClerkAPIService/WithdrawMoney"
	ClerkAPIService_CloseAccount_FullMethodName  = "/ClerkAPIService/CloseAccount"
)

// ClerkAPIServiceClient is the client API for ClerkAPIService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClerkAPIServiceClient interface {
	// Creates a new account and returns it
	OpenAccount(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*OpenAccountResponse, error)
	// Returns the list of open accounts
	ListAccounts(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListAccountsResponse, error)
	// Adds money to an account
	AddMoney(ctx context.Context, in *AddMoneyRequest, opts ...grpc.CallOption) (*AddMoneyResponse, error)
	// Removes money from an account
	WithdrawMoney(ctx context.Context, in *WithdrawMoneyRequest, opts ...grpc.CallOption) (*WithdrawMoneyResponse, error)
	// Close an account
	CloseAccount(ctx context.Context, in *CloseAccountRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type clerkAPIServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewClerkAPIServiceClient(cc grpc.ClientConnInterface) ClerkAPIServiceClient {
	return &clerkAPIServiceClient{cc}
}

func (c *clerkAPIServiceClient) OpenAccount(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*OpenAccountResponse, error) {
	out := new(OpenAccountResponse)
	err := c.cc.Invoke(ctx, ClerkAPIService_OpenAccount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clerkAPIServiceClient) ListAccounts(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListAccountsResponse, error) {
	out := new(ListAccountsResponse)
	err := c.cc.Invoke(ctx, ClerkAPIService_ListAccounts_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clerkAPIServiceClient) AddMoney(ctx context.Context, in *AddMoneyRequest, opts ...grpc.CallOption) (*AddMoneyResponse, error) {
	out := new(AddMoneyResponse)
	err := c.cc.Invoke(ctx, ClerkAPIService_AddMoney_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clerkAPIServiceClient) WithdrawMoney(ctx context.Context, in *WithdrawMoneyRequest, opts ...grpc.CallOption) (*WithdrawMoneyResponse, error) {
	out := new(WithdrawMoneyResponse)
	err := c.cc.Invoke(ctx, ClerkAPIService_WithdrawMoney_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clerkAPIServiceClient) CloseAccount(ctx context.Context, in *CloseAccountRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ClerkAPIService_CloseAccount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClerkAPIServiceServer is the server API for ClerkAPIService service.
// All implementations should embed UnimplementedClerkAPIServiceServer
// for forward compatibility
type ClerkAPIServiceServer interface {
	// Creates a new account and returns it
	OpenAccount(context.Context, *emptypb.Empty) (*OpenAccountResponse, error)
	// Returns the list of open accounts
	ListAccounts(context.Context, *emptypb.Empty) (*ListAccountsResponse, error)
	// Adds money to an account
	AddMoney(context.Context, *AddMoneyRequest) (*AddMoneyResponse, error)
	// Removes money from an account
	WithdrawMoney(context.Context, *WithdrawMoneyRequest) (*WithdrawMoneyResponse, error)
	// Close an account
	CloseAccount(context.Context, *CloseAccountRequest) (*emptypb.Empty, error)
}

// UnimplementedClerkAPIServiceServer should be embedded to have forward compatible implementations.
type UnimplementedClerkAPIServiceServer struct {
}

func (UnimplementedClerkAPIServiceServer) OpenAccount(context.Context, *emptypb.Empty) (*OpenAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpenAccount not implemented")
}
func (UnimplementedClerkAPIServiceServer) ListAccounts(context.Context, *emptypb.Empty) (*ListAccountsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAccounts not implemented")
}
func (UnimplementedClerkAPIServiceServer) AddMoney(context.Context, *AddMoneyRequest) (*AddMoneyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMoney not implemented")
}
func (UnimplementedClerkAPIServiceServer) WithdrawMoney(context.Context, *WithdrawMoneyRequest) (*WithdrawMoneyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WithdrawMoney not implemented")
}
func (UnimplementedClerkAPIServiceServer) CloseAccount(context.Context, *CloseAccountRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CloseAccount not implemented")
}

// UnsafeClerkAPIServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClerkAPIServiceServer will
// result in compilation errors.
type UnsafeClerkAPIServiceServer interface {
	mustEmbedUnimplementedClerkAPIServiceServer()
}

func RegisterClerkAPIServiceServer(s grpc.ServiceRegistrar, srv ClerkAPIServiceServer) {
	s.RegisterService(&ClerkAPIService_ServiceDesc, srv)
}

func _ClerkAPIService_OpenAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClerkAPIServiceServer).OpenAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClerkAPIService_OpenAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClerkAPIServiceServer).OpenAccount(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClerkAPIService_ListAccounts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClerkAPIServiceServer).ListAccounts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClerkAPIService_ListAccounts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClerkAPIServiceServer).ListAccounts(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClerkAPIService_AddMoney_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddMoneyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClerkAPIServiceServer).AddMoney(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClerkAPIService_AddMoney_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClerkAPIServiceServer).AddMoney(ctx, req.(*AddMoneyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClerkAPIService_WithdrawMoney_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WithdrawMoneyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClerkAPIServiceServer).WithdrawMoney(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClerkAPIService_WithdrawMoney_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClerkAPIServiceServer).WithdrawMoney(ctx, req.(*WithdrawMoneyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClerkAPIService_CloseAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloseAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClerkAPIServiceServer).CloseAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClerkAPIService_CloseAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClerkAPIServiceServer).CloseAccount(ctx, req.(*CloseAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ClerkAPIService_ServiceDesc is the grpc.ServiceDesc for ClerkAPIService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClerkAPIService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ClerkAPIService",
	HandlerType: (*ClerkAPIServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OpenAccount",
			Handler:    _ClerkAPIService_OpenAccount_Handler,
		},
		{
			MethodName: "ListAccounts",
			Handler:    _ClerkAPIService_ListAccounts_Handler,
		},
		{
			MethodName: "AddMoney",
			Handler:    _ClerkAPIService_AddMoney_Handler,
		},
		{
			MethodName: "WithdrawMoney",
			Handler:    _ClerkAPIService_WithdrawMoney_Handler,
		},
		{
			MethodName: "CloseAccount",
			Handler:    _ClerkAPIService_CloseAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
