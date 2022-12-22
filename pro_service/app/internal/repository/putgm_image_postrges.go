package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/jmoiron/sqlx"
)

type PutgmImageRepo struct {
	db *sqlx.DB
}

func NewPutgmImageRepo(db *sqlx.DB) *PutgmImageRepo {
	return &PutgmImageRepo{db: db}
}

func (r *PutgmImageRepo) Get(req *pro_api.GetPutgmImageRequest) (images []*pro_api.PutgmImage, err error) {
	query := fmt.Sprintf("SELECT id, form, gasket, url FROM %s WHERE form=$1", PUTGmImageTable)

	if err = r.db.Select(&images, query, req.Form); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return images, nil
}

func (r *PutgmImageRepo) Create(image *pro_api.CreatePutgmImageRequest) (id string, err error) {
	query := fmt.Sprintf(`INSERT INTO %s (form, gasket, url) VALUES ($1, $2, $3)  RETURNING id`, PUTGmImageTable)

	row := r.db.QueryRow(query, image.Form, image.Gasket, image.Url)

	var idInt int
	if err := row.Scan(&idInt); err != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *PutgmImageRepo) Update(image *pro_api.UpdatePutgmImageRequest) error {
	query := fmt.Sprintf("UPDATE %s SET form=$1, gasket=$2, url=$3 WHERE id=$4", PUTGmImageTable)

	id, err := strconv.Atoi(image.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, image.Form, image.Gasket, image.Url, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgmImageRepo) Delete(image *pro_api.DeletePutgmImageRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", PUTGmImageTable)

	id, err := strconv.Atoi(image.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
