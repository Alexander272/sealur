package repository

import (
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type Stand interface {
	GetAll(stand proto.GetStands) ([]proto.Stand, error)
	Create(stand proto.CreateStand) (string, error)
	Update(stand proto.UpdateStand) error
	Delete(stand proto.DeleteStand) error
}

type Repositories struct {
	Stand
}

func NewRepo(db *sqlx.DB) *Repositories {
	return &Repositories{
		Stand: NewStandRepo(db),
	}
}
