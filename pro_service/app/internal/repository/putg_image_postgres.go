package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type PutgImageRepo struct {
	db *sqlx.DB
}

func NewPutgImageRepo(db *sqlx.DB) *PutgImageRepo {
	return &PutgImageRepo{db: db}
}

func (r *PutgImageRepo) Get(req *proto.GetPutgImageRequest) (images []*proto.PutgImage, err error) {
	query := fmt.Sprintf("SELECT id, form, gasket, url FROM %s WHERE form=$1", PUTGImageTable)

	if err = r.db.Select(&images, query, req.Form); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return images, nil
}

func (r *PutgImageRepo) Create(image *proto.CreatePutgImageRequest) (id string, err error) {
	query := fmt.Sprintf(`INSERT INTO %s (form, gasket, url) VALUES ($1, $2, $3)  RETURNING id`, PUTGImageTable)

	row := r.db.QueryRow(query, image.Form, image.Gasket, image.Url)

	var idInt int
	if err := row.Scan(&idInt); err != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *PutgImageRepo) Update(image *proto.UpdatePutgImageRequest) error {
	query := fmt.Sprintf("UPDATE %s SET form=$1, gasket=$2, url=$3 WHERE id=$4", PUTGImageTable)

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

func (r *PutgImageRepo) Delete(image *proto.DeletePutgImageRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", PUTGImageTable)

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
