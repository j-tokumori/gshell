// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.17.3
// source: api/auth.proto

package api

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

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthServiceClient interface {
	// ユーザー認証
	Register(ctx context.Context, in *RegisterArgs, opts ...grpc.CallOption) (*RegisterReply, error)
	// ログイン
	Login(ctx context.Context, in *LoginArgs, opts ...grpc.CallOption) (*LoginReply, error)
	// データ引き継ぎURL取得
	GetTransferURL(ctx context.Context, in *GetTransferURLArgs, opts ...grpc.CallOption) (*GetTransferURLReply, error)
	// データ引き継ぎ実行
	ExecuteTransfer(ctx context.Context, in *ExecuteTransferArgs, opts ...grpc.CallOption) (*ExecuteTransferReply, error)
	// 非ログインデータ引き継ぎ連携URL取得
	GetConnectURLAndToken(ctx context.Context, in *GetConnectURLAndTokenArgs, opts ...grpc.CallOption) (*GetConnectURLAndTokenReply, error)
	// 非ログインデータ引き継ぎ連携実行実装
	ExecuteConnectProviderAndToken(ctx context.Context, in *ExecuteConnectProviderAndTokenArgs, opts ...grpc.CallOption) (*ExecuteConnectProviderAndTokenReply, error)
	// 非ログイン他ゲームユーザーのデータ引き継ぎ連携解除実装
	// 非ログイン状態でプロバイダ情報を元に、そのプロバイダと連携している別ゲームユーザーとの連携を解除する
	ConnectReleaseOtherAndToken(ctx context.Context, in *ConnectReleaseOtherAndTokenArgs, opts ...grpc.CallOption) (*ConnectReleaseOtherAndTokenReply, error)
	// 非ログインデータ引き継ぎ連携解除実装
	// 非ログイン状態でデータ引き継ぎのプロバイダとの連携を解除する
	ConnectReleaseAndToken(ctx context.Context, in *ConnectReleaseAndTokenArgs, opts ...grpc.CallOption) (*ConnectReleaseAndTokenReply, error)
	// データ引き継ぎが有効かどうかを返却
	GetTransferEnable(ctx context.Context, in *GetTransferEnableArgs, opts ...grpc.CallOption) (*GetTransferEnableReply, error)
	// サーバー死活フラグ取得
	// サーバーが現在有効(障害・メンテ中などでない)かどうか返却する
	GetServerEnable(ctx context.Context, in *GetServerEnableArgs, opts ...grpc.CallOption) (*GetServerEnableReply, error)
	// 非ログイン用データ引き継ぎ連携時のトークン検証
	VerifyConnectTokenAndUser(ctx context.Context, in *VerifyConnectTokenAndUserArgs, opts ...grpc.CallOption) (*VerifyConnectTokenAndUserReply, error)
	// 引き継ぎ時の検証
	VerifyTransferToken(ctx context.Context, in *VerifyTransferTokenArgs, opts ...grpc.CallOption) (*VerifyTransferTokenReply, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) Register(ctx context.Context, in *RegisterArgs, opts ...grpc.CallOption) (*RegisterReply, error) {
	out := new(RegisterReply)
	err := c.cc.Invoke(ctx, "/api.AuthService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) Login(ctx context.Context, in *LoginArgs, opts ...grpc.CallOption) (*LoginReply, error) {
	out := new(LoginReply)
	err := c.cc.Invoke(ctx, "/api.AuthService/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) GetTransferURL(ctx context.Context, in *GetTransferURLArgs, opts ...grpc.CallOption) (*GetTransferURLReply, error) {
	out := new(GetTransferURLReply)
	err := c.cc.Invoke(ctx, "/api.AuthService/GetTransferURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ExecuteTransfer(ctx context.Context, in *ExecuteTransferArgs, opts ...grpc.CallOption) (*ExecuteTransferReply, error) {
	out := new(ExecuteTransferReply)
	err := c.cc.Invoke(ctx, "/api.AuthService/ExecuteTransfer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) GetConnectURLAndToken(ctx context.Context, in *GetConnectURLAndTokenArgs, opts ...grpc.CallOption) (*GetConnectURLAndTokenReply, error) {
	out := new(GetConnectURLAndTokenReply)
	err := c.cc.Invoke(ctx, "/api.AuthService/GetConnectURLAndToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ExecuteConnectProviderAndToken(ctx context.Context, in *ExecuteConnectProviderAndTokenArgs, opts ...grpc.CallOption) (*ExecuteConnectProviderAndTokenReply, error) {
	out := new(ExecuteConnectProviderAndTokenReply)
	err := c.cc.Invoke(ctx, "/api.AuthService/ExecuteConnectProviderAndToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ConnectReleaseOtherAndToken(ctx context.Context, in *ConnectReleaseOtherAndTokenArgs, opts ...grpc.CallOption) (*ConnectReleaseOtherAndTokenReply, error) {
	out := new(ConnectReleaseOtherAndTokenReply)
	err := c.cc.Invoke(ctx, "/api.AuthService/ConnectReleaseOtherAndToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ConnectReleaseAndToken(ctx context.Context, in *ConnectReleaseAndTokenArgs, opts ...grpc.CallOption) (*ConnectReleaseAndTokenReply, error) {
	out := new(ConnectReleaseAndTokenReply)
	err := c.cc.Invoke(ctx, "/api.AuthService/ConnectReleaseAndToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) GetTransferEnable(ctx context.Context, in *GetTransferEnableArgs, opts ...grpc.CallOption) (*GetTransferEnableReply, error) {
	out := new(GetTransferEnableReply)
	err := c.cc.Invoke(ctx, "/api.AuthService/GetTransferEnable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) GetServerEnable(ctx context.Context, in *GetServerEnableArgs, opts ...grpc.CallOption) (*GetServerEnableReply, error) {
	out := new(GetServerEnableReply)
	err := c.cc.Invoke(ctx, "/api.AuthService/GetServerEnable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) VerifyConnectTokenAndUser(ctx context.Context, in *VerifyConnectTokenAndUserArgs, opts ...grpc.CallOption) (*VerifyConnectTokenAndUserReply, error) {
	out := new(VerifyConnectTokenAndUserReply)
	err := c.cc.Invoke(ctx, "/api.AuthService/VerifyConnectTokenAndUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) VerifyTransferToken(ctx context.Context, in *VerifyTransferTokenArgs, opts ...grpc.CallOption) (*VerifyTransferTokenReply, error) {
	out := new(VerifyTransferTokenReply)
	err := c.cc.Invoke(ctx, "/api.AuthService/VerifyTransferToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServiceServer is the server API for AuthService service.
// All implementations must embed UnimplementedAuthServiceServer
// for forward compatibility
type AuthServiceServer interface {
	// ユーザー認証
	Register(context.Context, *RegisterArgs) (*RegisterReply, error)
	// ログイン
	Login(context.Context, *LoginArgs) (*LoginReply, error)
	// データ引き継ぎURL取得
	GetTransferURL(context.Context, *GetTransferURLArgs) (*GetTransferURLReply, error)
	// データ引き継ぎ実行
	ExecuteTransfer(context.Context, *ExecuteTransferArgs) (*ExecuteTransferReply, error)
	// 非ログインデータ引き継ぎ連携URL取得
	GetConnectURLAndToken(context.Context, *GetConnectURLAndTokenArgs) (*GetConnectURLAndTokenReply, error)
	// 非ログインデータ引き継ぎ連携実行実装
	ExecuteConnectProviderAndToken(context.Context, *ExecuteConnectProviderAndTokenArgs) (*ExecuteConnectProviderAndTokenReply, error)
	// 非ログイン他ゲームユーザーのデータ引き継ぎ連携解除実装
	// 非ログイン状態でプロバイダ情報を元に、そのプロバイダと連携している別ゲームユーザーとの連携を解除する
	ConnectReleaseOtherAndToken(context.Context, *ConnectReleaseOtherAndTokenArgs) (*ConnectReleaseOtherAndTokenReply, error)
	// 非ログインデータ引き継ぎ連携解除実装
	// 非ログイン状態でデータ引き継ぎのプロバイダとの連携を解除する
	ConnectReleaseAndToken(context.Context, *ConnectReleaseAndTokenArgs) (*ConnectReleaseAndTokenReply, error)
	// データ引き継ぎが有効かどうかを返却
	GetTransferEnable(context.Context, *GetTransferEnableArgs) (*GetTransferEnableReply, error)
	// サーバー死活フラグ取得
	// サーバーが現在有効(障害・メンテ中などでない)かどうか返却する
	GetServerEnable(context.Context, *GetServerEnableArgs) (*GetServerEnableReply, error)
	// 非ログイン用データ引き継ぎ連携時のトークン検証
	VerifyConnectTokenAndUser(context.Context, *VerifyConnectTokenAndUserArgs) (*VerifyConnectTokenAndUserReply, error)
	// 引き継ぎ時の検証
	VerifyTransferToken(context.Context, *VerifyTransferTokenArgs) (*VerifyTransferTokenReply, error)
	mustEmbedUnimplementedAuthServiceServer()
}

// UnimplementedAuthServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServiceServer struct {
}

func (UnimplementedAuthServiceServer) Register(context.Context, *RegisterArgs) (*RegisterReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedAuthServiceServer) Login(context.Context, *LoginArgs) (*LoginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAuthServiceServer) GetTransferURL(context.Context, *GetTransferURLArgs) (*GetTransferURLReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransferURL not implemented")
}
func (UnimplementedAuthServiceServer) ExecuteTransfer(context.Context, *ExecuteTransferArgs) (*ExecuteTransferReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteTransfer not implemented")
}
func (UnimplementedAuthServiceServer) GetConnectURLAndToken(context.Context, *GetConnectURLAndTokenArgs) (*GetConnectURLAndTokenReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConnectURLAndToken not implemented")
}
func (UnimplementedAuthServiceServer) ExecuteConnectProviderAndToken(context.Context, *ExecuteConnectProviderAndTokenArgs) (*ExecuteConnectProviderAndTokenReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteConnectProviderAndToken not implemented")
}
func (UnimplementedAuthServiceServer) ConnectReleaseOtherAndToken(context.Context, *ConnectReleaseOtherAndTokenArgs) (*ConnectReleaseOtherAndTokenReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnectReleaseOtherAndToken not implemented")
}
func (UnimplementedAuthServiceServer) ConnectReleaseAndToken(context.Context, *ConnectReleaseAndTokenArgs) (*ConnectReleaseAndTokenReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnectReleaseAndToken not implemented")
}
func (UnimplementedAuthServiceServer) GetTransferEnable(context.Context, *GetTransferEnableArgs) (*GetTransferEnableReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransferEnable not implemented")
}
func (UnimplementedAuthServiceServer) GetServerEnable(context.Context, *GetServerEnableArgs) (*GetServerEnableReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetServerEnable not implemented")
}
func (UnimplementedAuthServiceServer) VerifyConnectTokenAndUser(context.Context, *VerifyConnectTokenAndUserArgs) (*VerifyConnectTokenAndUserReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyConnectTokenAndUser not implemented")
}
func (UnimplementedAuthServiceServer) VerifyTransferToken(context.Context, *VerifyTransferTokenArgs) (*VerifyTransferTokenReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyTransferToken not implemented")
}
func (UnimplementedAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {}

// UnsafeAuthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServiceServer will
// result in compilation errors.
type UnsafeAuthServiceServer interface {
	mustEmbedUnimplementedAuthServiceServer()
}

func RegisterAuthServiceServer(s grpc.ServiceRegistrar, srv AuthServiceServer) {
	s.RegisterService(&AuthService_ServiceDesc, srv)
}

func _AuthService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.AuthService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Register(ctx, req.(*RegisterArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.AuthService/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Login(ctx, req.(*LoginArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_GetTransferURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransferURLArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).GetTransferURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.AuthService/GetTransferURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).GetTransferURL(ctx, req.(*GetTransferURLArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ExecuteTransfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteTransferArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ExecuteTransfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.AuthService/ExecuteTransfer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ExecuteTransfer(ctx, req.(*ExecuteTransferArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_GetConnectURLAndToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConnectURLAndTokenArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).GetConnectURLAndToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.AuthService/GetConnectURLAndToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).GetConnectURLAndToken(ctx, req.(*GetConnectURLAndTokenArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ExecuteConnectProviderAndToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteConnectProviderAndTokenArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ExecuteConnectProviderAndToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.AuthService/ExecuteConnectProviderAndToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ExecuteConnectProviderAndToken(ctx, req.(*ExecuteConnectProviderAndTokenArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ConnectReleaseOtherAndToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectReleaseOtherAndTokenArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ConnectReleaseOtherAndToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.AuthService/ConnectReleaseOtherAndToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ConnectReleaseOtherAndToken(ctx, req.(*ConnectReleaseOtherAndTokenArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ConnectReleaseAndToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectReleaseAndTokenArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ConnectReleaseAndToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.AuthService/ConnectReleaseAndToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ConnectReleaseAndToken(ctx, req.(*ConnectReleaseAndTokenArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_GetTransferEnable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransferEnableArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).GetTransferEnable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.AuthService/GetTransferEnable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).GetTransferEnable(ctx, req.(*GetTransferEnableArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_GetServerEnable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetServerEnableArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).GetServerEnable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.AuthService/GetServerEnable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).GetServerEnable(ctx, req.(*GetServerEnableArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_VerifyConnectTokenAndUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyConnectTokenAndUserArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).VerifyConnectTokenAndUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.AuthService/VerifyConnectTokenAndUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).VerifyConnectTokenAndUser(ctx, req.(*VerifyConnectTokenAndUserArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_VerifyTransferToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyTransferTokenArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).VerifyTransferToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.AuthService/VerifyTransferToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).VerifyTransferToken(ctx, req.(*VerifyTransferTokenArgs))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthService_ServiceDesc is the grpc.ServiceDesc for AuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _AuthService_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _AuthService_Login_Handler,
		},
		{
			MethodName: "GetTransferURL",
			Handler:    _AuthService_GetTransferURL_Handler,
		},
		{
			MethodName: "ExecuteTransfer",
			Handler:    _AuthService_ExecuteTransfer_Handler,
		},
		{
			MethodName: "GetConnectURLAndToken",
			Handler:    _AuthService_GetConnectURLAndToken_Handler,
		},
		{
			MethodName: "ExecuteConnectProviderAndToken",
			Handler:    _AuthService_ExecuteConnectProviderAndToken_Handler,
		},
		{
			MethodName: "ConnectReleaseOtherAndToken",
			Handler:    _AuthService_ConnectReleaseOtherAndToken_Handler,
		},
		{
			MethodName: "ConnectReleaseAndToken",
			Handler:    _AuthService_ConnectReleaseAndToken_Handler,
		},
		{
			MethodName: "GetTransferEnable",
			Handler:    _AuthService_GetTransferEnable_Handler,
		},
		{
			MethodName: "GetServerEnable",
			Handler:    _AuthService_GetServerEnable_Handler,
		},
		{
			MethodName: "VerifyConnectTokenAndUser",
			Handler:    _AuthService_VerifyConnectTokenAndUser_Handler,
		},
		{
			MethodName: "VerifyTransferToken",
			Handler:    _AuthService_VerifyTransferToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/auth.proto",
}
