package service

import (
	"context"

	"github.com/Alexander272/sealur/user_service/internal/repo"
	proto_user "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto"
	proto_email "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto/email"
	"github.com/Alexander272/sealur/user_service/pkg/hasher"
)

type User interface {
	Get(context.Context, *proto_user.GetUserRequest) (*proto_user.User, error)
	GetAll(context.Context, *proto_user.GetAllUserRequest) ([]*proto_user.User, error)
	GetNewUsers(context.Context, *proto_user.GetNewUserRequest) ([]*proto_user.User, error)
	Create(context.Context, *proto_user.CreateUserRequest) (*proto_user.SuccessResponse, error)
	Confirm(context.Context, *proto_user.ConfirmUserRequest) (*proto_user.SuccessResponse, error)
	Update(context.Context, *proto_user.UpdateUserRequest) error
	Delete(context.Context, *proto_user.DeleteUserRequest) error
}

type Role interface {
	Get(context.Context, *proto_user.GetRolesRequest) ([]*proto_user.Role, error)
	Create(context.Context, *proto_user.CreateRoleRequest) (*proto_user.SuccessResponse, error)
	Update(context.Context, *proto_user.UpdateRoleRequest) error
	Delete(context.Context, *proto_user.DeleteRoleRequest) error
}

type Services struct {
	User
	Role
}

type Deps struct {
	Repos  *repo.Repo
	Email  proto_email.EmailServiceClient
	Hasher hasher.PasswordHasher
}

func NewServices(deps Deps) *Services {
	return &Services{
		User: NewUserService(deps.Repos.Users, deps.Repos.Role, deps.Hasher),
		Role: NewRoleService(deps.Repos.Role),
	}
}
