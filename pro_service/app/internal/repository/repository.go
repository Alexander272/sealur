package repository

import (
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type Stand interface {
	GetAll(stand *proto.GetStandsRequest) ([]*proto.Stand, error)
	Create(stand *proto.CreateStandRequest) (string, error)
	Update(stand *proto.UpdateStandRequest) error
	Delete(stand *proto.DeleteStandRequest) error
}

type Repositories struct {
	Stand
}

func NewRepo(db *sqlx.DB) *Repositories {
	return &Repositories{
		Stand: NewStandRepo(db),
	}
}
