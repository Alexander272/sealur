package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/user_service/internal/models"
	"github.com/Alexander272/sealur/user_service/internal/repo"
	proto_user "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/user_service/pkg/hasher"
)

type UserService struct {
	hasher   hasher.PasswordHasher
	userRepo repo.Users
	roleRepo repo.Role
}

func NewUserService(user repo.Users, role repo.Role, hasher hasher.PasswordHasher) *UserService {
	return &UserService{
		userRepo: user,
		roleRepo: role,
		hasher:   hasher,
	}
}

func (s *UserService) Get(ctx context.Context, req *proto_user.GetUserRequest) (u *proto_user.User, err error) {
	user, err := s.userRepo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user. error: %w", err)
	}

	if (user == models.User{}) {
		return nil, fmt.Errorf("unable to find user")
	}

	if req.Login != "" {
		salt := strings.Split(user.Password, ".")[1]
		pass, err := s.hasher.Hash(req.Password, salt)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password. error: %w", err)
		}

		if fmt.Sprintf("%s.%s", pass, salt) != user.Password {
			return nil, fmt.Errorf("passwords do not match")
		}
	}

	roles, err := s.roleRepo.Get(ctx, &proto_user.GetRolesRequest{UserId: user.Id})
	if err != nil {
		return nil, fmt.Errorf("failed to get roles. error: %w", err)
	}

	var userRoles []*proto_user.Role
	for _, r := range roles {
		ur := proto_user.Role(r)
		userRoles = append(userRoles, &ur)
	}

	u = &proto_user.User{
		Id:           user.Id,
		Organization: user.Organization,
		Name:         user.Name,
		Email:        user.Email,
		City:         user.City,
		Position:     user.Position,
		Phone:        user.Phone,
		Roles:        userRoles,
	}

	return u, nil
}
