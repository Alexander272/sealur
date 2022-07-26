package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/jmoiron/sqlx"
)

type AdditRepo struct {
	db *sqlx.DB
}

func NewAdditRepo(db *sqlx.DB) *AdditRepo {
	return &AdditRepo{db: db}
}

func (r *AdditRepo) GetAll() (addit []models.Addit, err error) {
	query := fmt.Sprintf(`SELECT id, materials, mod, temperature, mounting, graphite, fillers, coating, 
		construction, obturator, basis, p_obturator, sealant FROM %s LIMIT 1`, AdditionalTable)

	if err = r.db.Select(&addit, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return addit, nil
}

func (r *AdditRepo) Create(add *pro_api.CreateAddRequest) error {
	query := fmt.Sprintf(`INSERT INTO %s (materials, mod, temperature, mounting, graphite, fillers, coating, construction, obturator, basis, 
		p_obturator, sealant) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`, AdditionalTable)

	_, err := r.db.Exec(query, add.Materials, add.Mod, add.Temperature, add.Mounting, add.Graphite, add.Fillers, add.Coating, add.Construction,
		add.Obturator, add.Basis, add.PObturator, add.Sealant)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdateMat(mat models.UpdateMat) error {
	query := fmt.Sprintf("UPDATE %s SET materials=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(mat.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, mat.Materials, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdateMod(mod models.UpdateMod) error {
	query := fmt.Sprintf("UPDATE %s SET mod=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(mod.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, mod.Mod, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdateTemp(temp models.UpdateTemp) error {
	query := fmt.Sprintf("UPDATE %s SET temperature=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(temp.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, temp.Temperature, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdateMoun(moun models.UpdateMoun) error {
	query := fmt.Sprintf("UPDATE %s SET mounting=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(moun.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error %w", err)
	}

	_, err = r.db.Exec(query, moun.Mounting, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdateGrap(grap models.UpdateGrap) error {
	query := fmt.Sprintf("UPDATE %s SET graphite=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(grap.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, grap.Graphite, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdateFillers(fillers models.UpdateFill) error {
	query := fmt.Sprintf("UPDATE %s SET fillers=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(fillers.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, fillers.Fillers, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdateCoating(coating models.UpdateCoating) error {
	query := fmt.Sprintf("UPDATE %s SET coating=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(coating.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, coating.Coating, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdateConstruction(constr models.UpdateConstr) error {
	query := fmt.Sprintf("UPDATE %s SET construction=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(constr.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, constr.Construction, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdateObturator(obturator models.UpdateObturator) error {
	query := fmt.Sprintf("UPDATE %s SET obturator=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(obturator.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, obturator.Obturator, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdateBasis(basis models.UpdateBasis) error {
	query := fmt.Sprintf("UPDATE %s SET basis=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(basis.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, basis.Basis, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdatePObturator(pObt models.UpdatePObturator) error {
	query := fmt.Sprintf("UPDATE %s SET p_obturator=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(pObt.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, pObt.PObturator, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdateSealant(sealant models.UpdateSealant) error {
	query := fmt.Sprintf("UPDATE %s SET sealant=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(sealant.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, sealant.Sealant, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
