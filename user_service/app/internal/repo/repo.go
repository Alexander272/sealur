package repo

import (
	"context"

	"github.com/Alexander272/sealur/user_service/internal/config"
	"github.com/Alexander272/sealur/user_service/internal/models"
	proto_user "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type Users interface {
	Get(context.Context, *proto_user.GetUserRequest) (models.User, error)
	GetAll(context.Context, *proto_user.GetAllUserRequest) ([]models.User, error)
	GetNew(context.Context, *proto_user.GetNewUserRequest) ([]models.User, error)
	Create(context.Context, *proto_user.CreateUserRequest) error
	Confirm(context.Context, *proto_user.ConfirmUserRequest) (models.ConfirmUser, error)
	Update(context.Context, *proto_user.UpdateUserRequest) error
	Delete(context.Context, *proto_user.DeleteUserRequest) (models.DeleteUser, error)
}

type Role interface {
	Get(context.Context, *proto_user.GetRolesRequest) ([]models.Role, error)
	GetAll(context.Context, *proto_user.GetAllRolesRequest) ([]models.Role, error)
	Create(context.Context, []*proto_user.CreateRoleRequest) error
	Update(context.Context, *proto_user.UpdateRoleRequest) error
	Delete(context.Context, *proto_user.DeleteRoleRequest) error
}

type IP interface {
	GetAll(context.Context, *proto_user.GetAllIpRequest) ([]models.Ip, error)
	Add(context.Context, *proto_user.AddIpRequest) error
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
