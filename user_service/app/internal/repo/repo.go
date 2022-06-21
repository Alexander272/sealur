package repo

import (
	"context"

	"github.com/Alexander272/sealur/user_service/internal/models"
	proto_user "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type Users interface {
	Get(context.Context, *proto_user.GetUserRequest) (models.User, error)
	GetAll(context.Context, *proto_user.GetAllUserRequest) ([]models.User, error)
	GetNew(context.Context, *proto_user.GetNewUserRequest) ([]models.User, error)
	Create(context.Context, *proto_user.CreateUserRequest) error
	Confirm(context.Context, *proto_user.ConfirmUserRequest) (string, error)
	Update(context.Context, *proto_user.UpdateUserRequest) error
	Delete(context.Context, *proto_user.DeleteUserRequest) error
}

type Role interface {
	Get(context.Context, *proto_user.GetRolesRequest) ([]models.Role, error)
	GetAll(context.Context, *proto_user.GetAllRolesRequest) ([]models.Role, error)
	Create(context.Context, []*proto_user.CreateRoleRequest) error
	Update(context.Context, *proto_user.UpdateRoleRequest) error
	Delete(context.Context, *proto_user.DeleteRoleRequest) error
}

type Repo struct {
	Users
	Role
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		Users: NewUserRepo(db, UsersTable),
		Role:  NewRoleRepo(db, RolesTable),
	}
}
