package grpc

import (
	"context"

	proto_user "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto"
)

func (h *Handler) GetRoles(ctx context.Context, req *proto_user.GetRolesRequest) (*proto_user.RolesResponse, error) {
	roles, err := h.service.Role.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return &proto_user.RolesResponse{Roles: roles}, nil
}

func (h *Handler) CreateRole(ctx context.Context, role *proto_user.CreateRoleRequest) (*proto_user.SuccessResponse, error) {
	success, err := h.service.Role.Create(ctx, role)
	if err != nil {
		return nil, err
	}

	return success, nil
}

func (h *Handler) UpdateRole(ctx context.Context, role *proto_user.UpdateRoleRequest) (*proto_user.IdResponse, error) {
	if err := h.service.Role.Update(ctx, role); err != nil {
		return nil, err
	}

	return &proto_user.IdResponse{Id: role.Id}, nil
}

func (h *Handler) DeleteRole(ctx context.Context, role *proto_user.DeleteRoleRequest) (*proto_user.IdResponse, error) {
	if err := h.service.Role.Delete(ctx, role); err != nil {
		return nil, err
	}

	return &proto_user.IdResponse{Id: role.Id}, nil
}
