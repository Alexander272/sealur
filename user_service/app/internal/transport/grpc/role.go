package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/user_api"
)

func (h *Handler) GetRoles(ctx context.Context, req *user_api.GetRolesRequest) (*user_api.RolesResponse, error) {
	roles, err := h.service.Role.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return &user_api.RolesResponse{Roles: roles}, nil
}

func (h *Handler) CreateRole(ctx context.Context, role *user_api.CreateRoleRequest) (*user_api.SuccessResponse, error) {
	success, err := h.service.Role.Create(ctx, role)
	if err != nil {
		return nil, err
	}

	return success, nil
}

func (h *Handler) UpdateRole(ctx context.Context, role *user_api.UpdateRoleRequest) (*user_api.IdResponse, error) {
	if err := h.service.Role.Update(ctx, role); err != nil {
		return nil, err
	}

	return &user_api.IdResponse{Id: role.Id}, nil
}

func (h *Handler) DeleteRole(ctx context.Context, role *user_api.DeleteRoleRequest) (*user_api.IdResponse, error) {
	if err := h.service.Role.Delete(ctx, role); err != nil {
		return nil, err
	}

	return &user_api.IdResponse{Id: role.Id}, nil
}
