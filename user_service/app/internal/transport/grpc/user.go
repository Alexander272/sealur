package grpc

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Alexander272/sealur/user_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/user_api"
)

func (h *Handler) GetUser(ctx context.Context, req *user_api.GetUserRequest) (*user_api.UserResponse, error) {
	user, err := h.service.User.Get(ctx, req)
	if err != nil {
		if errors.Is(err, models.ErrPassword) || errors.Is(err, sql.ErrNoRows) {
			return &user_api.UserResponse{User: nil}, nil
		}
		return nil, err
	}

	return &user_api.UserResponse{User: user}, nil
}

func (h *Handler) GetAllUsers(ctx context.Context, req *user_api.GetAllUserRequest) (*user_api.UsersResponse, error) {
	users, err := h.service.User.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}

	return &user_api.UsersResponse{Users: users}, nil
}

func (h *Handler) GetNewUsers(ctx context.Context, req *user_api.GetNewUserRequest) (*user_api.UsersResponse, error) {
	users, err := h.service.User.GetNew(ctx, req)
	if err != nil {
		return nil, err
	}

	return &user_api.UsersResponse{Users: users}, nil
}

func (h *Handler) ConfirmUser(ctx context.Context, user *user_api.ConfirmUserRequest) (*user_api.SuccessResponse, error) {
	success, err := h.service.User.Confirm(ctx, user)
	if err != nil {
		return nil, err
	}

	return success, err
}

func (h *Handler) CreateUser(ctx context.Context, user *user_api.CreateUserRequest) (*user_api.SuccessResponse, error) {
	success, err := h.service.User.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return success, err
}

func (h *Handler) UpdateUser(ctx context.Context, user *user_api.UpdateUserRequest) (*user_api.IdResponse, error) {
	if err := h.service.User.Update(ctx, user); err != nil {
		return nil, err
	}

	return &user_api.IdResponse{Id: user.Id}, nil
}

func (h *Handler) DeleteUser(ctx context.Context, user *user_api.DeleteUserRequest) (*user_api.IdResponse, error) {
	if err := h.service.User.Delete(ctx, user); err != nil {
		return nil, err
	}

	return &user_api.IdResponse{Id: user.Id}, nil
}

func (h *Handler) RejectUser(ctx context.Context, user *user_api.DeleteUserRequest) (*user_api.IdResponse, error) {
	if err := h.service.User.Reject(ctx, user); err != nil {
		return nil, err
	}

	return &user_api.IdResponse{Id: user.Id}, nil
}
