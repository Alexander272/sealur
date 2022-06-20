// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/user.proto

/*
Package proto_user is a generated protocol buffer package.

It is generated from these files:
	proto/user.proto

It has these top-level messages:
	IdResponse
	SuccessResponse
	PingRequest
	PingResponse
	GetUserRequest
	GetAllUserRequest
	GetNewUserRequest
	CreateUserRequest
	UpdateUserRequest
	DeleteUserRequest
	ConfirmUserRequest
	User
	UserResponse
	UsersResponse
	GetRolesRequest
	GetAllRolesRequest
	CreateRoleRequest
	UpdateRoleRequest
	DeleteRoleRequest
	Role
	RoleResponse
*/
package proto_user

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type IdResponse struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *IdResponse) Reset()                    { *m = IdResponse{} }
func (m *IdResponse) String() string            { return proto.CompactTextString(m) }
func (*IdResponse) ProtoMessage()               {}
func (*IdResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *IdResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type SuccessResponse struct {
	Success bool `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
}

func (m *SuccessResponse) Reset()                    { *m = SuccessResponse{} }
func (m *SuccessResponse) String() string            { return proto.CompactTextString(m) }
func (*SuccessResponse) ProtoMessage()               {}
func (*SuccessResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SuccessResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type PingRequest struct {
}

func (m *PingRequest) Reset()                    { *m = PingRequest{} }
func (m *PingRequest) String() string            { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()               {}
func (*PingRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type PingResponse struct {
	Ping string `protobuf:"bytes,1,opt,name=ping" json:"ping,omitempty"`
}

func (m *PingResponse) Reset()                    { *m = PingResponse{} }
func (m *PingResponse) String() string            { return proto.CompactTextString(m) }
func (*PingResponse) ProtoMessage()               {}
func (*PingResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *PingResponse) GetPing() string {
	if m != nil {
		return m.Ping
	}
	return ""
}

// User Service ----------------------------------------------------------------------
type GetUserRequest struct {
	Login    string `protobuf:"bytes,1,opt,name=login" json:"login,omitempty"`
	Id       string `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
	Password string `protobuf:"bytes,3,opt,name=password" json:"password,omitempty"`
}

func (m *GetUserRequest) Reset()                    { *m = GetUserRequest{} }
func (m *GetUserRequest) String() string            { return proto.CompactTextString(m) }
func (*GetUserRequest) ProtoMessage()               {}
func (*GetUserRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GetUserRequest) GetLogin() string {
	if m != nil {
		return m.Login
	}
	return ""
}

func (m *GetUserRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *GetUserRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type GetAllUserRequest struct {
}

func (m *GetAllUserRequest) Reset()                    { *m = GetAllUserRequest{} }
func (m *GetAllUserRequest) String() string            { return proto.CompactTextString(m) }
func (*GetAllUserRequest) ProtoMessage()               {}
func (*GetAllUserRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type GetNewUserRequest struct {
}

func (m *GetNewUserRequest) Reset()                    { *m = GetNewUserRequest{} }
func (m *GetNewUserRequest) String() string            { return proto.CompactTextString(m) }
func (*GetNewUserRequest) ProtoMessage()               {}
func (*GetNewUserRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type CreateUserRequest struct {
	Organization string `protobuf:"bytes,1,opt,name=organization" json:"organization,omitempty"`
	Name         string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Email        string `protobuf:"bytes,3,opt,name=email" json:"email,omitempty"`
	City         string `protobuf:"bytes,4,opt,name=city" json:"city,omitempty"`
	Position     string `protobuf:"bytes,5,opt,name=position" json:"position,omitempty"`
	Phone        string `protobuf:"bytes,6,opt,name=phone" json:"phone,omitempty"`
}

func (m *CreateUserRequest) Reset()                    { *m = CreateUserRequest{} }
func (m *CreateUserRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateUserRequest) ProtoMessage()               {}
func (*CreateUserRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *CreateUserRequest) GetOrganization() string {
	if m != nil {
		return m.Organization
	}
	return ""
}

func (m *CreateUserRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateUserRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *CreateUserRequest) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *CreateUserRequest) GetPosition() string {
	if m != nil {
		return m.Position
	}
	return ""
}

func (m *CreateUserRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

type UpdateUserRequest struct {
	Id       string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name     string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Email    string `protobuf:"bytes,3,opt,name=email" json:"email,omitempty"`
	Position string `protobuf:"bytes,4,opt,name=position" json:"position,omitempty"`
	Phone    string `protobuf:"bytes,5,opt,name=phone" json:"phone,omitempty"`
}

func (m *UpdateUserRequest) Reset()                    { *m = UpdateUserRequest{} }
func (m *UpdateUserRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateUserRequest) ProtoMessage()               {}
func (*UpdateUserRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *UpdateUserRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *UpdateUserRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *UpdateUserRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *UpdateUserRequest) GetPosition() string {
	if m != nil {
		return m.Position
	}
	return ""
}

func (m *UpdateUserRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

type DeleteUserRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *DeleteUserRequest) Reset()                    { *m = DeleteUserRequest{} }
func (m *DeleteUserRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteUserRequest) ProtoMessage()               {}
func (*DeleteUserRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *DeleteUserRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type ConfirmUserRequest struct {
	Id       string  `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Login    string  `protobuf:"bytes,2,opt,name=login" json:"login,omitempty"`
	Password string  `protobuf:"bytes,3,opt,name=password" json:"password,omitempty"`
	Roles    []*Role `protobuf:"bytes,4,rep,name=roles" json:"roles,omitempty"`
}

func (m *ConfirmUserRequest) Reset()                    { *m = ConfirmUserRequest{} }
func (m *ConfirmUserRequest) String() string            { return proto.CompactTextString(m) }
func (*ConfirmUserRequest) ProtoMessage()               {}
func (*ConfirmUserRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *ConfirmUserRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ConfirmUserRequest) GetLogin() string {
	if m != nil {
		return m.Login
	}
	return ""
}

func (m *ConfirmUserRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *ConfirmUserRequest) GetRoles() []*Role {
	if m != nil {
		return m.Roles
	}
	return nil
}

type User struct {
	Id           string  `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Organization string  `protobuf:"bytes,2,opt,name=organization" json:"organization,omitempty"`
	Name         string  `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
	Email        string  `protobuf:"bytes,4,opt,name=email" json:"email,omitempty"`
	City         string  `protobuf:"bytes,5,opt,name=city" json:"city,omitempty"`
	Position     string  `protobuf:"bytes,6,opt,name=position" json:"position,omitempty"`
	Phone        string  `protobuf:"bytes,7,opt,name=phone" json:"phone,omitempty"`
	Roles        []*Role `protobuf:"bytes,8,rep,name=roles" json:"roles,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *User) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *User) GetOrganization() string {
	if m != nil {
		return m.Organization
	}
	return ""
}

func (m *User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *User) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *User) GetPosition() string {
	if m != nil {
		return m.Position
	}
	return ""
}

func (m *User) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *User) GetRoles() []*Role {
	if m != nil {
		return m.Roles
	}
	return nil
}

type UserResponse struct {
	User *User `protobuf:"bytes,1,opt,name=user" json:"user,omitempty"`
}

func (m *UserResponse) Reset()                    { *m = UserResponse{} }
func (m *UserResponse) String() string            { return proto.CompactTextString(m) }
func (*UserResponse) ProtoMessage()               {}
func (*UserResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *UserResponse) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type UsersResponse struct {
	User []*User `protobuf:"bytes,1,rep,name=user" json:"user,omitempty"`
}

func (m *UsersResponse) Reset()                    { *m = UsersResponse{} }
func (m *UsersResponse) String() string            { return proto.CompactTextString(m) }
func (*UsersResponse) ProtoMessage()               {}
func (*UsersResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *UsersResponse) GetUser() []*User {
	if m != nil {
		return m.User
	}
	return nil
}

// Role Service ----------------------------------------------------------------------
type GetRolesRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=userId" json:"userId,omitempty"`
}

func (m *GetRolesRequest) Reset()                    { *m = GetRolesRequest{} }
func (m *GetRolesRequest) String() string            { return proto.CompactTextString(m) }
func (*GetRolesRequest) ProtoMessage()               {}
func (*GetRolesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *GetRolesRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

type GetAllRolesRequest struct {
}

func (m *GetAllRolesRequest) Reset()                    { *m = GetAllRolesRequest{} }
func (m *GetAllRolesRequest) String() string            { return proto.CompactTextString(m) }
func (*GetAllRolesRequest) ProtoMessage()               {}
func (*GetAllRolesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

type CreateRoleRequest struct {
	UserId  string `protobuf:"bytes,1,opt,name=userId" json:"userId,omitempty"`
	Service string `protobuf:"bytes,2,opt,name=service" json:"service,omitempty"`
	Role    string `protobuf:"bytes,3,opt,name=role" json:"role,omitempty"`
}

func (m *CreateRoleRequest) Reset()                    { *m = CreateRoleRequest{} }
func (m *CreateRoleRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateRoleRequest) ProtoMessage()               {}
func (*CreateRoleRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *CreateRoleRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *CreateRoleRequest) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *CreateRoleRequest) GetRole() string {
	if m != nil {
		return m.Role
	}
	return ""
}

type UpdateRoleRequest struct {
	Id      string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Service string `protobuf:"bytes,2,opt,name=service" json:"service,omitempty"`
	Role    string `protobuf:"bytes,3,opt,name=role" json:"role,omitempty"`
}

func (m *UpdateRoleRequest) Reset()                    { *m = UpdateRoleRequest{} }
func (m *UpdateRoleRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateRoleRequest) ProtoMessage()               {}
func (*UpdateRoleRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

func (m *UpdateRoleRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *UpdateRoleRequest) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *UpdateRoleRequest) GetRole() string {
	if m != nil {
		return m.Role
	}
	return ""
}

type DeleteRoleRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *DeleteRoleRequest) Reset()                    { *m = DeleteRoleRequest{} }
func (m *DeleteRoleRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteRoleRequest) ProtoMessage()               {}
func (*DeleteRoleRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{18} }

func (m *DeleteRoleRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type Role struct {
	Id      string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Service string `protobuf:"bytes,2,opt,name=service" json:"service,omitempty"`
	Role    string `protobuf:"bytes,3,opt,name=role" json:"role,omitempty"`
}

func (m *Role) Reset()                    { *m = Role{} }
func (m *Role) String() string            { return proto.CompactTextString(m) }
func (*Role) ProtoMessage()               {}
func (*Role) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{19} }

func (m *Role) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Role) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *Role) GetRole() string {
	if m != nil {
		return m.Role
	}
	return ""
}

type RoleResponse struct {
	Role []*Role `protobuf:"bytes,1,rep,name=role" json:"role,omitempty"`
}

func (m *RoleResponse) Reset()                    { *m = RoleResponse{} }
func (m *RoleResponse) String() string            { return proto.CompactTextString(m) }
func (*RoleResponse) ProtoMessage()               {}
func (*RoleResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{20} }

func (m *RoleResponse) GetRole() []*Role {
	if m != nil {
		return m.Role
	}
	return nil
}

func init() {
	proto.RegisterType((*IdResponse)(nil), "proto_user.IdResponse")
	proto.RegisterType((*SuccessResponse)(nil), "proto_user.SuccessResponse")
	proto.RegisterType((*PingRequest)(nil), "proto_user.PingRequest")
	proto.RegisterType((*PingResponse)(nil), "proto_user.PingResponse")
	proto.RegisterType((*GetUserRequest)(nil), "proto_user.GetUserRequest")
	proto.RegisterType((*GetAllUserRequest)(nil), "proto_user.GetAllUserRequest")
	proto.RegisterType((*GetNewUserRequest)(nil), "proto_user.GetNewUserRequest")
	proto.RegisterType((*CreateUserRequest)(nil), "proto_user.CreateUserRequest")
	proto.RegisterType((*UpdateUserRequest)(nil), "proto_user.UpdateUserRequest")
	proto.RegisterType((*DeleteUserRequest)(nil), "proto_user.DeleteUserRequest")
	proto.RegisterType((*ConfirmUserRequest)(nil), "proto_user.ConfirmUserRequest")
	proto.RegisterType((*User)(nil), "proto_user.User")
	proto.RegisterType((*UserResponse)(nil), "proto_user.UserResponse")
	proto.RegisterType((*UsersResponse)(nil), "proto_user.UsersResponse")
	proto.RegisterType((*GetRolesRequest)(nil), "proto_user.GetRolesRequest")
	proto.RegisterType((*GetAllRolesRequest)(nil), "proto_user.GetAllRolesRequest")
	proto.RegisterType((*CreateRoleRequest)(nil), "proto_user.CreateRoleRequest")
	proto.RegisterType((*UpdateRoleRequest)(nil), "proto_user.UpdateRoleRequest")
	proto.RegisterType((*DeleteRoleRequest)(nil), "proto_user.DeleteRoleRequest")
	proto.RegisterType((*Role)(nil), "proto_user.Role")
	proto.RegisterType((*RoleResponse)(nil), "proto_user.RoleResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for UserService service

type UserServiceClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	// User Service ----------------------------------------------------------------------
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*UserResponse, error)
	GetAllUsers(ctx context.Context, in *GetAllUserRequest, opts ...grpc.CallOption) (*UsersResponse, error)
	GetNewUsers(ctx context.Context, in *GetNewUserRequest, opts ...grpc.CallOption) (*UsersResponse, error)
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*SuccessResponse, error)
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*IdResponse, error)
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*IdResponse, error)
	ConfirmUser(ctx context.Context, in *ConfirmUserRequest, opts ...grpc.CallOption) (*SuccessResponse, error)
}

type userServiceClient struct {
	cc *grpc.ClientConn
}

func NewUserServiceClient(cc *grpc.ClientConn) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := grpc.Invoke(ctx, "/proto_user.UserService/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := grpc.Invoke(ctx, "/proto_user.UserService/GetUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetAllUsers(ctx context.Context, in *GetAllUserRequest, opts ...grpc.CallOption) (*UsersResponse, error) {
	out := new(UsersResponse)
	err := grpc.Invoke(ctx, "/proto_user.UserService/GetAllUsers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetNewUsers(ctx context.Context, in *GetNewUserRequest, opts ...grpc.CallOption) (*UsersResponse, error) {
	out := new(UsersResponse)
	err := grpc.Invoke(ctx, "/proto_user.UserService/GetNewUsers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*SuccessResponse, error) {
	out := new(SuccessResponse)
	err := grpc.Invoke(ctx, "/proto_user.UserService/CreateUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*IdResponse, error) {
	out := new(IdResponse)
	err := grpc.Invoke(ctx, "/proto_user.UserService/UpdateUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*IdResponse, error) {
	out := new(IdResponse)
	err := grpc.Invoke(ctx, "/proto_user.UserService/DeleteUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) ConfirmUser(ctx context.Context, in *ConfirmUserRequest, opts ...grpc.CallOption) (*SuccessResponse, error) {
	out := new(SuccessResponse)
	err := grpc.Invoke(ctx, "/proto_user.UserService/ConfirmUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserService service

type UserServiceServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	// User Service ----------------------------------------------------------------------
	GetUser(context.Context, *GetUserRequest) (*UserResponse, error)
	GetAllUsers(context.Context, *GetAllUserRequest) (*UsersResponse, error)
	GetNewUsers(context.Context, *GetNewUserRequest) (*UsersResponse, error)
	CreateUser(context.Context, *CreateUserRequest) (*SuccessResponse, error)
	UpdateUser(context.Context, *UpdateUserRequest) (*IdResponse, error)
	DeleteUser(context.Context, *DeleteUserRequest) (*IdResponse, error)
	ConfirmUser(context.Context, *ConfirmUserRequest) (*SuccessResponse, error)
}

func RegisterUserServiceServer(s *grpc.Server, srv UserServiceServer) {
	s.RegisterService(&_UserService_serviceDesc, srv)
}

func _UserService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_user.UserService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_user.UserService/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetAllUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetAllUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_user.UserService/GetAllUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetAllUsers(ctx, req.(*GetAllUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetNewUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNewUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetNewUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_user.UserService/GetNewUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetNewUsers(ctx, req.(*GetNewUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_user.UserService/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_user.UserService/UpdateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_user.UserService/DeleteUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ConfirmUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfirmUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ConfirmUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_user.UserService/ConfirmUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ConfirmUser(ctx, req.(*ConfirmUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto_user.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _UserService_Ping_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _UserService_GetUser_Handler,
		},
		{
			MethodName: "GetAllUsers",
			Handler:    _UserService_GetAllUsers_Handler,
		},
		{
			MethodName: "GetNewUsers",
			Handler:    _UserService_GetNewUsers_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _UserService_CreateUser_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _UserService_UpdateUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _UserService_DeleteUser_Handler,
		},
		{
			MethodName: "ConfirmUser",
			Handler:    _UserService_ConfirmUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/user.proto",
}

func init() { proto.RegisterFile("proto/user.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 651 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x55, 0xdd, 0x6e, 0xd3, 0x30,
	0x14, 0x56, 0xdb, 0xf4, 0x87, 0xd3, 0x76, 0x5d, 0xcd, 0x34, 0x42, 0xf8, 0xd1, 0x64, 0x10, 0x1a,
	0x42, 0xea, 0xa4, 0x31, 0x2e, 0x10, 0x17, 0x08, 0x3a, 0xa9, 0x8c, 0x0b, 0x04, 0x99, 0x76, 0x01,
	0x37, 0x28, 0xb4, 0x87, 0x62, 0x29, 0x8d, 0x43, 0x9c, 0x32, 0xc1, 0x05, 0x3c, 0x0e, 0x0f, 0xc1,
	0x5b, 0xf0, 0x44, 0xc8, 0x8e, 0xd3, 0x38, 0x49, 0x13, 0x69, 0xe2, 0x2a, 0x3e, 0x3e, 0xc7, 0x9f,
	0xbf, 0xe3, 0xcf, 0xfe, 0x02, 0xbb, 0x61, 0xc4, 0x63, 0x7e, 0xb4, 0x16, 0x18, 0x4d, 0xd4, 0x90,
	0x80, 0xfa, 0x7c, 0x94, 0x33, 0xf4, 0x36, 0xc0, 0xd9, 0xc2, 0x45, 0x11, 0xf2, 0x40, 0x20, 0xd9,
	0x81, 0x26, 0x5b, 0xd8, 0x8d, 0x83, 0xc6, 0xe1, 0x35, 0xb7, 0xc9, 0x16, 0xf4, 0x11, 0x8c, 0xce,
	0xd7, 0xf3, 0x39, 0x0a, 0xb1, 0x29, 0xb1, 0xa1, 0x2b, 0x92, 0x29, 0x55, 0xd7, 0x73, 0xd3, 0x90,
	0x0e, 0xa1, 0xff, 0x96, 0x05, 0x4b, 0x17, 0xbf, 0xae, 0x51, 0xc4, 0x94, 0xc2, 0x20, 0x09, 0xf5,
	0x42, 0x02, 0x56, 0xc8, 0x82, 0xa5, 0x46, 0x57, 0x63, 0xea, 0xc2, 0xce, 0x0c, 0xe3, 0x0b, 0x81,
	0x91, 0x5e, 0x45, 0xf6, 0xa0, 0xed, 0xf3, 0x25, 0x0b, 0x74, 0x59, 0x12, 0x68, 0x5e, 0xcd, 0x94,
	0x17, 0x71, 0xa0, 0x17, 0x7a, 0x42, 0x5c, 0xf2, 0x68, 0x61, 0xb7, 0xd4, 0xec, 0x26, 0xa6, 0xd7,
	0x61, 0x3c, 0xc3, 0xf8, 0x85, 0xef, 0x1b, 0xb0, 0x7a, 0xf2, 0x0d, 0x5e, 0x9a, 0x93, 0xbf, 0x1b,
	0x30, 0x9e, 0x46, 0xe8, 0xc5, 0x68, 0x32, 0xa0, 0x30, 0xe0, 0xd1, 0xd2, 0x0b, 0xd8, 0x0f, 0x2f,
	0x66, 0x3c, 0x25, 0x92, 0x9b, 0x93, 0xbd, 0x04, 0xde, 0x0a, 0x35, 0x23, 0x35, 0x96, 0xcc, 0x71,
	0xe5, 0x31, 0x5f, 0x13, 0x4a, 0x02, 0x59, 0x39, 0x67, 0xf1, 0x77, 0xdb, 0x4a, 0x2a, 0xe5, 0x58,
	0xb1, 0xe7, 0x82, 0x29, 0xf4, 0xb6, 0x66, 0xaf, 0x63, 0x89, 0x12, 0x7e, 0xe1, 0x01, 0xda, 0x9d,
	0x04, 0x45, 0x05, 0xf4, 0x17, 0x8c, 0x2f, 0xc2, 0x45, 0x81, 0x68, 0x41, 0xac, 0x2b, 0x90, 0x32,
	0x09, 0x58, 0x55, 0x04, 0xda, 0x26, 0x81, 0x7b, 0x30, 0x3e, 0x45, 0x1f, 0x6b, 0x09, 0xd0, 0x9f,
	0x40, 0xa6, 0x3c, 0xf8, 0xcc, 0xa2, 0x55, 0x1d, 0xcd, 0x8d, 0xc2, 0x4d, 0x53, 0xe1, 0x1a, 0x45,
	0xc9, 0x03, 0x68, 0x47, 0xdc, 0x47, 0x61, 0x5b, 0x07, 0xad, 0xc3, 0xfe, 0xf1, 0xee, 0x24, 0xbb,
	0xbf, 0x13, 0x97, 0xfb, 0xe8, 0x26, 0x69, 0xfa, 0xb7, 0x01, 0x96, 0xdc, 0xb9, 0xb4, 0x65, 0x51,
	0xd2, 0x66, 0x8d, 0xa4, 0xad, 0x6d, 0xa7, 0x67, 0x6d, 0x93, 0xb4, 0x5d, 0x21, 0x69, 0xa7, 0xea,
	0x44, 0xbb, 0xc6, 0x89, 0x66, 0x4d, 0xf5, 0xea, 0x9b, 0x3a, 0x81, 0x41, 0x72, 0x9a, 0xfa, 0x19,
	0xdd, 0x07, 0x4b, 0xd6, 0xa8, 0xee, 0x0a, 0xcb, 0x54, 0x9d, 0xca, 0xd2, 0x27, 0x30, 0x94, 0x91,
	0xd8, 0xb2, 0xac, 0x55, 0xb3, 0xec, 0x21, 0x8c, 0x66, 0x18, 0xcb, 0xed, 0x45, 0x2a, 0xdf, 0x3e,
	0x74, 0x64, 0xea, 0x2c, 0x3d, 0x4f, 0x1d, 0xd1, 0x3d, 0x20, 0xc9, 0x33, 0x33, 0xab, 0xe9, 0xfb,
	0xf4, 0x45, 0xa9, 0x16, 0xea, 0x21, 0x94, 0x95, 0x60, 0xf4, 0x8d, 0xcd, 0xd3, 0x3b, 0x9b, 0x86,
	0xf2, 0x88, 0x65, 0xf7, 0xa9, 0x18, 0x72, 0x4c, 0xdf, 0xa5, 0x6f, 0xc0, 0x84, 0x2e, 0x2a, 0x7d,
	0x35, 0xc8, 0xcd, 0xad, 0xae, 0x81, 0xa4, 0xa7, 0x60, 0xc9, 0xf4, 0x7f, 0x6e, 0x75, 0x02, 0x83,
	0x64, 0x93, 0x4c, 0x0f, 0x55, 0xd3, 0xa8, 0x50, 0x5f, 0x65, 0x8f, 0xff, 0x58, 0xd0, 0x97, 0xf2,
	0x9c, 0x6b, 0xe4, 0xa7, 0x60, 0x49, 0x4f, 0x25, 0x37, 0xcc, 0x7a, 0xc3, 0x74, 0x1d, 0xbb, 0x9c,
	0xd0, 0x1b, 0x3e, 0x87, 0xae, 0xb6, 0x5a, 0xe2, 0x98, 0x45, 0x79, 0xff, 0xcd, 0x03, 0xe4, 0x2e,
	0xde, 0x0c, 0xfa, 0x99, 0xaf, 0x0a, 0x72, 0xa7, 0x00, 0x92, 0x37, 0x5c, 0xe7, 0x66, 0x11, 0x47,
	0x14, 0x80, 0xb4, 0x17, 0x97, 0x81, 0xf2, 0x26, 0x5d, 0x07, 0xf4, 0x0a, 0x20, 0xb3, 0xef, 0x3c,
	0x4e, 0xc9, 0xd6, 0x9d, 0x5b, 0x66, 0xba, 0xf8, 0x53, 0x9b, 0x02, 0x64, 0xfe, 0x9a, 0x47, 0x2a,
	0xf9, 0xae, 0xb3, 0x6f, 0xa6, 0x8d, 0x9f, 0xe7, 0x14, 0x20, 0xf3, 0xc8, 0x3c, 0x48, 0xc9, 0x3b,
	0x2b, 0x41, 0x5e, 0x43, 0xdf, 0xf0, 0x50, 0x72, 0x37, 0xd7, 0x54, 0xc9, 0x5c, 0x6b, 0xbb, 0x7a,
	0x39, 0xfa, 0x30, 0x9c, 0x1c, 0x3d, 0xcb, 0x0a, 0x3e, 0x75, 0xd4, 0xf8, 0xf1, 0xbf, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x98, 0xce, 0xcc, 0x86, 0x13, 0x08, 0x00, 0x00,
}