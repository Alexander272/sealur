package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/user_service/internal/repo"
	proto_user "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto"
)

type RoleService struct {
	repo repo.Role
}

func NewRoleService(repo repo.Role) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) Get(ctx context.Context, req *proto_user.GetRolesRequest) ([]*proto_user.Role, error) {
	roles, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles. error: %w", err)
	}

	var userRoles []*proto_user.Role
	for _, item := range roles {
		r := proto_user.Role{
			Id:      item.Id,
			Service: item.Service,
			Role:    item.Role,
		}
		userRoles = append(userRoles, &r)
	}

	return userRoles, nil
}

func (s *RoleService) Create(ctx context.Context, role *proto_user.CreateRoleRequest) (*proto_user.SuccessResponse, error) {
	if err := s.repo.Create(ctx, []*proto_user.CreateRoleRequest{role}); err != nil {
		return nil, fmt.Errorf("failed to create role. error: %w", err)
	}
	return &proto_user.SuccessResponse{Success: true}, nil
}

func (s *RoleService) Update(ctx context.Context, role *proto_user.UpdateRoleRequest) error {
	if err := s.repo.Update(ctx, role); err != nil {
		return fmt.Errorf("failed to update role. error: %w", err)
	}
	return nil
}

func (s *RoleService) Delete(ctx context.Context, role *proto_user.DeleteRoleRequest) error {
	if err := s.repo.Delete(ctx, role); err != nil {
		return fmt.Errorf("failed to delete role. error: %w", err)
	}
	return nil
}
