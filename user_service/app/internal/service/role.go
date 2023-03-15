package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/user_service/internal/repo"
	"github.com/Alexander272/sealur_proto/api/user/models/role_model"
)

type RoleService struct {
	repo repo.Role
}

func NewRoleService(repo repo.Role) *RoleService {
	return &RoleService{
		repo: repo,
	}
}

func (s *RoleService) GetDefault(ctx context.Context) (*role_model.Role, error) {
	role, err := s.repo.GetDefault(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get default role. error: %w", err)
	}
	return role, nil
}

// type RoleService struct {
// 	repo repo.Role
// }

// func NewRoleService(repo repo.Role) *RoleService {
// 	return &RoleService{repo: repo}
// }

// func (s *RoleService) Get(ctx context.Context, req *user_api.GetRolesRequest) ([]*user_api.Role, error) {
// 	roles, err := s.repo.Get(ctx, req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get roles. error: %w", err)
// 	}

// 	var userRoles []*user_api.Role
// 	for _, item := range roles {
// 		r := user_api.Role{
// 			Id:      item.Id,
// 			Service: item.Service,
// 			Role:    item.Role,
// 		}
// 		userRoles = append(userRoles, &r)
// 	}

// 	return userRoles, nil
// }

// func (s *RoleService) Create(ctx context.Context, role *user_api.CreateRoleRequest) (*user_api.SuccessResponse, error) {
// 	if err := s.repo.Create(ctx, []*user_api.CreateRoleRequest{role}); err != nil {
// 		return nil, fmt.Errorf("failed to create role. error: %w", err)
// 	}
// 	return &user_api.SuccessResponse{Success: true}, nil
// }

// func (s *RoleService) Update(ctx context.Context, role *user_api.UpdateRoleRequest) error {
// 	if err := s.repo.Update(ctx, role); err != nil {
// 		return fmt.Errorf("failed to update role. error: %w", err)
// 	}
// 	return nil
// }

// func (s *RoleService) Delete(ctx context.Context, role *user_api.DeleteRoleRequest) error {
// 	if err := s.repo.Delete(ctx, role); err != nil {
// 		return fmt.Errorf("failed to delete role. error: %w", err)
// 	}
// 	return nil
// }
