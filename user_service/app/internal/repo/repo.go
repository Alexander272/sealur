package repo

import (
	"context"

	"github.com/Alexander272/sealur/user_service/internal/config"
	"github.com/Alexander272/sealur/user_service/internal/repo/postgres"
	"github.com/Alexander272/sealur_proto/api/user/models/role_model"
	"github.com/Alexander272/sealur_proto/api/user/models/user_model"
	"github.com/Alexander272/sealur_proto/api/user/user_api"
	"github.com/jmoiron/sqlx"
)

type Users interface {
	Get(context.Context, *user_api.GetUser) (*user_model.User, error)
	GetFull(context.Context, *user_api.GetUser) (*user_model.FullUser, error)
	GetByEmail(context.Context, *user_api.GetUserByEmail) (*user_model.User, string, error)
	GetManagers(context.Context, *user_api.GetNewUser) ([]*user_model.User, error)
	GetAnalytics(context.Context, *user_api.GetUserAnalytics) (*user_api.Analytics, error)
	GetFullAnalytics(context.Context, *user_api.GetUsersByParam) ([]*user_model.AnalyticUsers, error)
	// Get(context.Context, *user_api.GetUser) (models.User, error)
	// GetAll(context.Context, *user_api.GetAllUser) ([]models.User, error)
	// GetNew(context.Context, *user_api.GetNewUser) ([]models.User, error)
	Create(context.Context, *user_api.CreateUser, string) (string, error)
	Confirm(context.Context, *user_api.ConfirmUser) error
	SetManager(context.Context, *user_api.UserManager) error
	Update(context.Context, *user_api.UpdateUser) error
	Visit(context.Context, *user_api.GetUser) error
	// Update(context.Context, *user_api.UpdateUser) error
	// Delete(context.Context, *user_api.DeleteUser) error
}

type Role interface {
	GetDefault(context.Context) (*role_model.Role, error)
}

type Region interface {
	GetManagerByRegion(ctx context.Context, region string) (*user_model.User, error)
}

// type Role interface {
// 	Get(context.Context, *user_api.GetRoles) ([]models.Role, error)
// 	GetAll(context.Context, *user_api.GetAllRoles) ([]models.Role, error)
// 	Create(context.Context, []*user_api.CreateRole) error
// 	Update(context.Context, *user_api.UpdateRole) error
// 	Delete(context.Context, *user_api.DeleteRole) error
// }

// type IP interface {
// 	GetAll(context.Context, *user_api.GetAllIp) ([]models.Ip, error)
// 	Add(context.Context, *user_api.AddIp) error
// }

type Repo struct {
	Users
	Role
	Region
	// IP
}

func NewRepo(db *sqlx.DB, conf *config.Config) *Repo {
	return &Repo{
		Users:  postgres.NewUserRepo(db),
		Role:   postgres.NewRoleRepo(db),
		Region: postgres.NewRegionRepo(db),
		// Users: postgres.NewUserRepo(db),
		// Role:  NewRoleRepo(db, RolesTable),
		// IP:    NewIpRepo(db, IpTable, conf.IP),
	}
}
