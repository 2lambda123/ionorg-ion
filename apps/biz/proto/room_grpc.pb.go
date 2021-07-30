// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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

// RoomClient is the client API for Room service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoomClient interface {
	Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinReply, error)
	Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*LeaveReply, error)
	GetParticipants(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetParticipantsReply, error)
	ReceiveNotification(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Room_ReceiveNotificationClient, error)
	SetImportance(ctx context.Context, in *SetImportanceRequest, opts ...grpc.CallOption) (*SetImportanceReply, error)
	LockConference(ctx context.Context, in *LockConferenceRequest, opts ...grpc.CallOption) (*LockConferenceReply, error)
	EndConference(ctx context.Context, in *EndConferenceRequest, opts ...grpc.CallOption) (*EndConferenceReply, error)
	EditParticipantInfo(ctx context.Context, in *EditParticipantInfoRequest, opts ...grpc.CallOption) (*EditParticipantInfoReply, error)
	AddParticipant(ctx context.Context, in *AddParticipantRequest, opts ...grpc.CallOption) (*AddParticipantReply, error)
	RemoveParticipant(ctx context.Context, in *RemoveParticipantRequest, opts ...grpc.CallOption) (*RemoveParticipantReply, error)
	SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageReply, error)
}

type roomClient struct {
	cc grpc.ClientConnInterface
}

func NewRoomClient(cc grpc.ClientConnInterface) RoomClient {
	return &roomClient{cc}
}

func (c *roomClient) Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinReply, error) {
	out := new(JoinReply)
	err := c.cc.Invoke(ctx, "/room.Room/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*LeaveReply, error) {
	out := new(LeaveReply)
	err := c.cc.Invoke(ctx, "/room.Room/Leave", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) GetParticipants(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetParticipantsReply, error) {
	out := new(GetParticipantsReply)
	err := c.cc.Invoke(ctx, "/room.Room/GetParticipants", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) ReceiveNotification(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Room_ReceiveNotificationClient, error) {
	stream, err := c.cc.NewStream(ctx, &Room_ServiceDesc.Streams[0], "/room.Room/ReceiveNotification", opts...)
	if err != nil {
		return nil, err
	}
	x := &roomReceiveNotificationClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Room_ReceiveNotificationClient interface {
	Recv() (*Notification, error)
	grpc.ClientStream
}

type roomReceiveNotificationClient struct {
	grpc.ClientStream
}

func (x *roomReceiveNotificationClient) Recv() (*Notification, error) {
	m := new(Notification)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *roomClient) SetImportance(ctx context.Context, in *SetImportanceRequest, opts ...grpc.CallOption) (*SetImportanceReply, error) {
	out := new(SetImportanceReply)
	err := c.cc.Invoke(ctx, "/room.Room/SetImportance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) LockConference(ctx context.Context, in *LockConferenceRequest, opts ...grpc.CallOption) (*LockConferenceReply, error) {
	out := new(LockConferenceReply)
	err := c.cc.Invoke(ctx, "/room.Room/LockConference", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) EndConference(ctx context.Context, in *EndConferenceRequest, opts ...grpc.CallOption) (*EndConferenceReply, error) {
	out := new(EndConferenceReply)
	err := c.cc.Invoke(ctx, "/room.Room/EndConference", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) EditParticipantInfo(ctx context.Context, in *EditParticipantInfoRequest, opts ...grpc.CallOption) (*EditParticipantInfoReply, error) {
	out := new(EditParticipantInfoReply)
	err := c.cc.Invoke(ctx, "/room.Room/EditParticipantInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) AddParticipant(ctx context.Context, in *AddParticipantRequest, opts ...grpc.CallOption) (*AddParticipantReply, error) {
	out := new(AddParticipantReply)
	err := c.cc.Invoke(ctx, "/room.Room/AddParticipant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) RemoveParticipant(ctx context.Context, in *RemoveParticipantRequest, opts ...grpc.CallOption) (*RemoveParticipantReply, error) {
	out := new(RemoveParticipantReply)
	err := c.cc.Invoke(ctx, "/room.Room/RemoveParticipant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageReply, error) {
	out := new(SendMessageReply)
	err := c.cc.Invoke(ctx, "/room.Room/SendMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoomServer is the server API for Room service.
// All implementations must embed UnimplementedRoomServer
// for forward compatibility
type RoomServer interface {
	Join(context.Context, *JoinRequest) (*JoinReply, error)
	Leave(context.Context, *LeaveRequest) (*LeaveReply, error)
	GetParticipants(context.Context, *Empty) (*GetParticipantsReply, error)
	ReceiveNotification(*Empty, Room_ReceiveNotificationServer) error
	SetImportance(context.Context, *SetImportanceRequest) (*SetImportanceReply, error)
	LockConference(context.Context, *LockConferenceRequest) (*LockConferenceReply, error)
	EndConference(context.Context, *EndConferenceRequest) (*EndConferenceReply, error)
	EditParticipantInfo(context.Context, *EditParticipantInfoRequest) (*EditParticipantInfoReply, error)
	AddParticipant(context.Context, *AddParticipantRequest) (*AddParticipantReply, error)
	RemoveParticipant(context.Context, *RemoveParticipantRequest) (*RemoveParticipantReply, error)
	SendMessage(context.Context, *SendMessageRequest) (*SendMessageReply, error)
	mustEmbedUnimplementedRoomServer()
}

// UnimplementedRoomServer must be embedded to have forward compatible implementations.
type UnimplementedRoomServer struct {
}

func (UnimplementedRoomServer) Join(context.Context, *JoinRequest) (*JoinReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedRoomServer) Leave(context.Context, *LeaveRequest) (*LeaveReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Leave not implemented")
}
func (UnimplementedRoomServer) GetParticipants(context.Context, *Empty) (*GetParticipantsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetParticipants not implemented")
}
func (UnimplementedRoomServer) ReceiveNotification(*Empty, Room_ReceiveNotificationServer) error {
	return status.Errorf(codes.Unimplemented, "method ReceiveNotification not implemented")
}
func (UnimplementedRoomServer) SetImportance(context.Context, *SetImportanceRequest) (*SetImportanceReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetImportance not implemented")
}
func (UnimplementedRoomServer) LockConference(context.Context, *LockConferenceRequest) (*LockConferenceReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LockConference not implemented")
}
func (UnimplementedRoomServer) EndConference(context.Context, *EndConferenceRequest) (*EndConferenceReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EndConference not implemented")
}
func (UnimplementedRoomServer) EditParticipantInfo(context.Context, *EditParticipantInfoRequest) (*EditParticipantInfoReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditParticipantInfo not implemented")
}
func (UnimplementedRoomServer) AddParticipant(context.Context, *AddParticipantRequest) (*AddParticipantReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddParticipant not implemented")
}
func (UnimplementedRoomServer) RemoveParticipant(context.Context, *RemoveParticipantRequest) (*RemoveParticipantReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveParticipant not implemented")
}
func (UnimplementedRoomServer) SendMessage(context.Context, *SendMessageRequest) (*SendMessageReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedRoomServer) mustEmbedUnimplementedRoomServer() {}

// UnsafeRoomServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoomServer will
// result in compilation errors.
type UnsafeRoomServer interface {
	mustEmbedUnimplementedRoomServer()
}

func RegisterRoomServer(s grpc.ServiceRegistrar, srv RoomServer) {
	s.RegisterService(&Room_ServiceDesc, srv)
}

func _Room_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/room.Room/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).Join(ctx, req.(*JoinRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_Leave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).Leave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/room.Room/Leave",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).Leave(ctx, req.(*LeaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_GetParticipants_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).GetParticipants(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/room.Room/GetParticipants",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).GetParticipants(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_ReceiveNotification_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RoomServer).ReceiveNotification(m, &roomReceiveNotificationServer{stream})
}

type Room_ReceiveNotificationServer interface {
	Send(*Notification) error
	grpc.ServerStream
}

type roomReceiveNotificationServer struct {
	grpc.ServerStream
}

func (x *roomReceiveNotificationServer) Send(m *Notification) error {
	return x.ServerStream.SendMsg(m)
}

func _Room_SetImportance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetImportanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).SetImportance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/room.Room/SetImportance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).SetImportance(ctx, req.(*SetImportanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_LockConference_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LockConferenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).LockConference(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/room.Room/LockConference",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).LockConference(ctx, req.(*LockConferenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_EndConference_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndConferenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).EndConference(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/room.Room/EndConference",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).EndConference(ctx, req.(*EndConferenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_EditParticipantInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditParticipantInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).EditParticipantInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/room.Room/EditParticipantInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).EditParticipantInfo(ctx, req.(*EditParticipantInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_AddParticipant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddParticipantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).AddParticipant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/room.Room/AddParticipant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).AddParticipant(ctx, req.(*AddParticipantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_RemoveParticipant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveParticipantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).RemoveParticipant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/room.Room/RemoveParticipant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).RemoveParticipant(ctx, req.(*RemoveParticipantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/room.Room/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).SendMessage(ctx, req.(*SendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Room_ServiceDesc is the grpc.ServiceDesc for Room service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Room_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "room.Room",
	HandlerType: (*RoomServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Join",
			Handler:    _Room_Join_Handler,
		},
		{
			MethodName: "Leave",
			Handler:    _Room_Leave_Handler,
		},
		{
			MethodName: "GetParticipants",
			Handler:    _Room_GetParticipants_Handler,
		},
		{
			MethodName: "SetImportance",
			Handler:    _Room_SetImportance_Handler,
		},
		{
			MethodName: "LockConference",
			Handler:    _Room_LockConference_Handler,
		},
		{
			MethodName: "EndConference",
			Handler:    _Room_EndConference_Handler,
		},
		{
			MethodName: "EditParticipantInfo",
			Handler:    _Room_EditParticipantInfo_Handler,
		},
		{
			MethodName: "AddParticipant",
			Handler:    _Room_AddParticipant_Handler,
		},
		{
			MethodName: "RemoveParticipant",
			Handler:    _Room_RemoveParticipant_Handler,
		},
		{
			MethodName: "SendMessage",
			Handler:    _Room_SendMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ReceiveNotification",
			Handler:       _Room_ReceiveNotification_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "apps/biz/proto/room.proto",
}
