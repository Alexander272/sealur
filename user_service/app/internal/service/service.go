package service

import (
	"context"

	"github.com/Alexander272/sealur/user_service/internal/repo"
	"github.com/Alexander272/sealur/user_service/pkg/hasher"
	"github.com/Alexander272/sealur_proto/api/email_api"
	"github.com/Alexander272/sealur_proto/api/user/models/role_model"
	"github.com/Alexander272/sealur_proto/api/user/user_api"
)

type User interface {
	Create(context.Context, *user_api.CreateUser) error
	// Get(context.Context, *user_api.GetUserRequest) (*user_api.User, error)
	// GetAll(context.Context, *user_api.GetAllUserRequest) ([]*user_api.User, int, error)
	// GetNew(context.Context, *user_api.GetNewUserRequest) ([]*user_api.User, error)
	// Create(context.Context, *user_api.CreateUserRequest) (*user_api.SuccessResponse, error)
	// Confirm(context.Context, *user_api.ConfirmUserRequest) (*user_api.SuccessResponse, error)
	// Update(context.Context, *user_api.UpdateUserRequest) error
	// Delete(context.Context, *user_api.DeleteUserRequest) error
	// Reject(ctx context.Context, user *user_api.DeleteUserRequest) error
}

type Role interface {
	GetDefault(context.Context) (*role_model.Role, error)
	// Get(context.Context, *user_api.GetRolesRequest) ([]*user_api.Role, error)
	// Create(context.Context, *user_api.CreateRoleRequest) (*user_api.SuccessResponse, error)
	// Update(context.Context, *user_api.UpdateRoleRequest) error
	// Delete(context.Context, *user_api.DeleteRoleRequest) error
}

type IP interface {
	// Add(ctx context.Context, ip *user_api.AddIpRequest) error
}

type Services struct {
	User
	Role
	IP
}

type Deps struct {
	Repos  *repo.Repo
	Email  email_api.EmailServiceClient
	Hasher hasher.PasswordHasher
}

func NewServices(deps Deps) *Services {
	return &Services{
		// User: NewUserService(deps.Repos.Users, deps.Repos.Role, deps.Repos.IP, deps.Hasher, deps.Email),
		// Role: NewRoleService(deps.Repos.Role),
		// IP:   NewIpService(deps.Repos.IP),
	}
}
