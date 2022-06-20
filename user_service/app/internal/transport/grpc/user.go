package grpc

import (
	"context"

	proto_user "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto"
)

func (h *Handler) GetUser(ctx context.Context, req *proto_user.GetUserRequest) (*proto_user.UserResponse, error) {
	return nil, nil
}

func (h *Handler) GetAllUsers(ctx context.Context, req *proto_user.GetAllUserRequest) (*proto_user.UsersResponse, error)

func (h *Handler) GetNewUsers(ctx context.Context, req *proto_user.GetNewUserRequest) (*proto_user.UsersResponse, error) {
	return nil, nil
}

func (h *Handler) ConfirmUser(ctx context.Context, user *proto_user.ConfirmUserRequest) (*proto_user.SuccessResponse, error)

func (h *Handler) CreateUser(ctx context.Context, user *proto_user.CreateUserRequest) (*proto_user.SuccessResponse, error)

func (h *Handler) UpdateUser(ctx context.Context, user *proto_user.UpdateUserRequest) (*proto_user.IdResponse, error)

func (h *Handler) DeleteUser(ctx context.Context, user *proto_user.DeleteUserRequest) (*proto_user.IdResponse, error)
