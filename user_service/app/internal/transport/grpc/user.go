package grpc

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Alexander272/sealur/user_service/internal/models"
	proto_user "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto"
)

func (h *Handler) GetUser(ctx context.Context, req *proto_user.GetUserRequest) (*proto_user.UserResponse, error) {
	user, err := h.service.User.Get(ctx, req)
	if err != nil {
		if errors.Is(err, models.ErrPassword) || errors.Is(err, sql.ErrNoRows) {
			return &proto_user.UserResponse{User: nil}, nil
		}
		return nil, err
	}

	return &proto_user.UserResponse{User: user}, nil
}

func (h *Handler) GetAllUsers(ctx context.Context, req *proto_user.GetAllUserRequest) (*proto_user.UsersResponse, error) {
	users, err := h.service.User.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}

	return &proto_user.UsersResponse{Users: users}, nil
}

func (h *Handler) GetNewUsers(ctx context.Context, req *proto_user.GetNewUserRequest) (*proto_user.UsersResponse, error) {
	users, err := h.service.User.GetNew(ctx, req)
	if err != nil {
		return nil, err
	}

	return &proto_user.UsersResponse{Users: users}, nil
}

func (h *Handler) ConfirmUser(ctx context.Context, user *proto_user.ConfirmUserRequest) (*proto_user.SuccessResponse, error) {
	success, err := h.service.User.Confirm(ctx, user)
	if err != nil {
		return nil, err
	}

	return success, err
}

func (h *Handler) CreateUser(ctx context.Context, user *proto_user.CreateUserRequest) (*proto_user.SuccessResponse, error) {
	success, err := h.service.User.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return success, err
}

func (h *Handler) UpdateUser(ctx context.Context, user *proto_user.UpdateUserRequest) (*proto_user.IdResponse, error) {
	if err := h.service.User.Update(ctx, user); err != nil {
		return nil, err
	}

	return &proto_user.IdResponse{Id: user.Id}, nil
}

func (h *Handler) DeleteUser(ctx context.Context, user *proto_user.DeleteUserRequest) (*proto_user.IdResponse, error) {
	if err := h.service.User.Delete(ctx, user); err != nil {
		return nil, err
	}

	return &proto_user.IdResponse{Id: user.Id}, nil
}

func (h *Handler) RejectUser(ctx context.Context, user *proto_user.DeleteUserRequest) (*proto_user.IdResponse, error) {
	if err := h.service.User.Reject(ctx, user); err != nil {
		return nil, err
	}

	return &proto_user.IdResponse{Id: user.Id}, nil
}
