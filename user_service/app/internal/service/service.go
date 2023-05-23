package service

import (
	"context"

	"github.com/Alexander272/sealur/user_service/internal/repo"
	"github.com/Alexander272/sealur/user_service/pkg/hasher"
	"github.com/Alexander272/sealur_proto/api/email_api"
	"github.com/Alexander272/sealur_proto/api/user/models/role_model"
	"github.com/Alexander272/sealur_proto/api/user/models/user_model"
	"github.com/Alexander272/sealur_proto/api/user/user_api"
)

type User interface {
	Get(context.Context, *user_api.GetUser) (*user_model.User, error)
	GetFull(context.Context, *user_api.GetUser) (*user_model.FullUser, error)
	GetByEmail(context.Context, *user_api.GetUserByEmail) (*user_model.User, error)
	GetManager(context.Context, *user_api.GetUser) (*user_api.Manager, error)
	GetManagers(context.Context, *user_api.GetNewUser) ([]*user_model.User, error)
	GetAnalytics(context.Context, *user_api.GetUserAnalytics) (*user_api.Analytics, error)
	GetFullAnalytics(context.Context, *user_api.GetUsersByParam) ([]*user_model.AnalyticUsers, error)
	Create(context.Context, *user_api.CreateUser) (string, error)
	Confirm(context.Context, *user_api.ConfirmUser) (*user_model.User, error)
	SetManager(context.Context, *user_api.UserManager) error
	Update(context.Context, *user_api.UpdateUser) error
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

type Region interface {
	GetManagerByRegion(ctx context.Context, region string) (*user_model.User, error)
}

type Services struct {
	User
	Role
	IP
	Region
}

type Deps struct {
	Repos  *repo.Repo
	Email  email_api.EmailServiceClient
	Hasher hasher.PasswordHasher
}

func NewServices(deps Deps) *Services {
	role := NewRoleService(deps.Repos.Role)
	region := NewRegionService(deps.Repos.Region)
	user := NewUserService(deps.Repos.Users, deps.Hasher, role, region)

	return &Services{
		Role:   role,
		User:   user,
		Region: region,
		// User: NewUserService(deps.Repos.Users, deps.Repos.Role, deps.Repos.IP, deps.Hasher, deps.Email),
		// Role: NewRoleService(deps.Repos.Role),
		// IP:   NewIpService(deps.Repos.IP),
	}
}
