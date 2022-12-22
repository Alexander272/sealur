package repo

import (
	"context"

	"github.com/Alexander272/sealur/user_service/internal/config"
	"github.com/Alexander272/sealur/user_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/user_api"
	"github.com/jmoiron/sqlx"
)

type Users interface {
	Get(context.Context, *user_api.GetUserRequest) (models.User, error)
	GetAll(context.Context, *user_api.GetAllUserRequest) ([]models.User, error)
	GetNew(context.Context, *user_api.GetNewUserRequest) ([]models.User, error)
	Create(context.Context, *user_api.CreateUserRequest) error
	Confirm(context.Context, *user_api.ConfirmUserRequest) (models.ConfirmUser, error)
	Update(context.Context, *user_api.UpdateUserRequest) error
	Delete(context.Context, *user_api.DeleteUserRequest) (models.DeleteUser, error)
}

type Role interface {
	Get(context.Context, *user_api.GetRolesRequest) ([]models.Role, error)
	GetAll(context.Context, *user_api.GetAllRolesRequest) ([]models.Role, error)
	Create(context.Context, []*user_api.CreateRoleRequest) error
	Update(context.Context, *user_api.UpdateRoleRequest) error
	Delete(context.Context, *user_api.DeleteRoleRequest) error
}

type IP interface {
	GetAll(context.Context, *user_api.GetAllIpRequest) ([]models.Ip, error)
	Add(context.Context, *user_api.AddIpRequest) error
}

type Repo struct {
	Users
	Role
	IP
}

func NewRepo(db *sqlx.DB, conf *config.Config) *Repo {
	return &Repo{
		Users: NewUserRepo(db, UsersTable),
		Role:  NewRoleRepo(db, RolesTable),
		IP:    NewIpRepo(db, IpTable, conf.IP),
	}
}
