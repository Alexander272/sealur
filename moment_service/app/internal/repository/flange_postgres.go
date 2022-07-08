package repository

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/jmoiron/sqlx"
)

type FlangeRepo struct {
	db *sqlx.DB
}

func NewFlangeRepo(db *sqlx.DB) *FlangeRepo {
	return &FlangeRepo{db: db}
}

var flangeSizes = []models.FlangeSize{
	{Pn: 1.0, D: 400, D1: 535, D2: 495, S0: (412 - 400) / 2, S1: (432 - 400) / 2, Lenght: 65 - 35, B: 35, Count: 20, Diameter: 20},
	{Pn: 1.6, D: 400, D1: 535, D2: 495, S0: (412 - 400) / 2, S1: (436 - 400) / 2, Lenght: 70 - 35, B: 35, Count: 20, Diameter: 20},
	{Pn: 2.5, D: 400, D1: 535, D2: 495, S0: (418 - 400) / 2, S1: (440 - 400) / 2, Lenght: 75 - 40, B: 40, Count: 24, Diameter: 20},
	{Pn: 4.0, D: 400, D1: 590, D2: 530, S0: (424 - 400) / 2, S1: (454 - 400) / 2, Lenght: 95 - 50, B: 50, Count: 20, Diameter: 30},
	{Pn: 6.3, D: 400, D1: 590, D2: 530, S0: (428 - 400) / 2, S1: (460 - 400) / 2, Lenght: 120 - 70, B: 70, Count: 20, Diameter: 30},
	{Pn: 1.0, D: 450, D1: 590, D2: 550, S0: (464 - 450) / 2, S1: (482 - 450) / 2, Lenght: 65 - 35, B: 35, Count: 24, Diameter: 20},
	{Pn: 1.6, D: 450, D1: 590, D2: 550, S0: (464 - 450) / 2, S1: (486 - 450) / 2, Lenght: 70 - 35, B: 35, Count: 24, Diameter: 20},
	{Pn: 2.5, D: 450, D1: 590, D2: 550, S0: (472 - 450) / 2, S1: (490 - 450) / 2, Lenght: 75 - 45, B: 45, Count: 24, Diameter: 20},
	{Pn: 4.0, D: 450, D1: 640, D2: 580, S0: (474 - 450) / 2, S1: (510 - 450) / 2, Lenght: 75 - 50, B: 50, Count: 20, Diameter: 30},
	{Pn: 6.3, D: 450, D1: 640, D2: 580, S0: (478 - 450) / 2, S1: (510 - 450) / 2, Lenght: 120 - 75, B: 75, Count: 20, Diameter: 30},
}

func (r *FlangeRepo) GetSize(ctx context.Context, dn, pn float64) (models.FlangeSize, error) {
	size := find(dn, pn)
	return size, nil
}

func find(dn, pn float64) models.FlangeSize {
	var res models.FlangeSize
	for _, fs := range flangeSizes {
		if fs.Pn == pn && fs.D == dn {
			res = fs
			break
		}
	}
	return res
}
