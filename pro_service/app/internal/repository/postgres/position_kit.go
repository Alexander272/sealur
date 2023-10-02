package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/rings_kit_construction_model"
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PositionKitRepo struct {
	db *sqlx.DB
}

func NewPositionKitRepo(db *sqlx.DB) *PositionKitRepo {
	return &PositionKitRepo{
		db: db,
	}
}

type PositionKit interface {
	Get(ctx context.Context, orderId string) (positions []*position_model.FullPosition, err error)
	GetFull(ctx context.Context, positionsId []string) ([]*position_model.OrderPositionRingsKit, error)
	Create(ctx context.Context, position *position_model.FullPosition) error
	Update(ctx context.Context, position *position_model.FullPosition) error
	Copy(ctx context.Context, targetPositionId string, position *position_api.CopyPosition) (string, error)
}

func (r *PositionKitRepo) Get(ctx context.Context, orderId string) (positions []*position_model.FullPosition, err error) {
	var data []models.KitPosition
	query := fmt.Sprintf(`SELECT p.id, title, amount, type, p.count, info,
		type_code, construction_code, r.count as rings_count, size, thickness, materials, modifying, drawing
		FROM %s AS p LEFT JOIN %s as r ON p.id=position_id
		WHERE order_id=$1 AND type=$2 ORDER BY count`,
		PositionTable, PositionRingsKitTable,
	)

	if err := r.db.Select(&data, query, orderId, position_model.PositionType_RingsKit.String()); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, rp := range data {
		positions = append(positions, &position_model.FullPosition{
			Id:     rp.Id,
			Count:  rp.Count,
			Title:  rp.Title,
			Amount: rp.Amount,
			Info:   rp.Info,
			Type:   position_model.PositionType_RingsKit,
			KitData: &position_model.PositionRingsKit{
				Type:             rp.TypeCode,
				ConstructionCode: rp.ConstructionCode,
				Count:            rp.RingsCount,
				Size:             rp.Size,
				Thickness:        rp.Thickness,
				Material:         rp.Material,
				Modifying:        rp.Modifying,
				Drawing:          rp.Drawing,
			},
		})
	}

	return positions, nil
}

func (r *PositionKitRepo) GetFull(ctx context.Context, positionsId []string) (positions []*position_model.OrderPositionRingsKit, err error) {
	var data []models.RingsKit
	query := fmt.Sprintf(`SELECT k.id, position_id, k.type_id, type_code, k.count, size, thickness, materials, modifying, drawing,
		construction_id, construction_code, title, image, same_rings, material_types, has_thickness, enabled_materials
		FROM %s AS k INNER JOIN %s AS c ON c.id=construction_id WHERE array[position_id] <@ $1 ORDER BY position_id`,
		PositionRingsKitTable, RingsKitConstructionTable,
	)

	if err := r.db.Select(&data, query, pq.Array(positionsId)); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, rk := range data {
		positions = append(positions, &position_model.OrderPositionRingsKit{
			Id:         rk.Id,
			PositionId: rk.PositionId,
			TypeId:     rk.TypeId,
			Type:       rk.TypeCode,
			Construction: &rings_kit_construction_model.RingsKitConstruction{
				Id:               rk.ConstructionId,
				Code:             rk.ConstructionCode,
				Title:            rk.Title,
				Image:            rk.Image,
				SameRings:        rk.SameRings,
				MaterialTypes:    rk.MaterialTypes,
				HasThickness:     rk.HasThickness,
				EnabledMaterials: rk.EnabledMaterials,
			},
			Count:     rk.Count,
			Size:      rk.Size,
			Thickness: rk.Thickness,
			Material:  rk.Material,
			Modifying: rk.Modifying,
			Drawing:   rk.Drawing,
		})
	}

	return positions, nil
}

func (r *PositionKitRepo) Create(ctx context.Context, position *position_model.FullPosition) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, position_id, type_id, type_code, construction_id, construction_code, count, size, thickness,
		materials, modifying, drawing) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		PositionRingsKitTable,
	)
	id := uuid.New()
	kit := position.KitData

	_, err := r.db.Exec(query, id, position.Id, kit.TypeId, kit.Type, kit.ConstructionId, kit.ConstructionCode, kit.Count, kit.Size,
		kit.Thickness, kit.Material, kit.Modifying, kit.Drawing,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *PositionKitRepo) Update(ctx context.Context, position *position_model.FullPosition) error {
	query := fmt.Sprintf(`UPDATE %s SET type_id=$1, type_code=$2, construction_id=$3, construction_code=$4, count=$5, size=$6, thickness=$7,
		materials=$8, modifying=$9, drawing=$10 WHERE position_id=$11`,
		PositionRingsKitTable,
	)

	kit := position.KitData
	_, err := r.db.Exec(query, kit.TypeId, kit.Type, kit.ConstructionId, kit.ConstructionCode, kit.Count, kit.Size, kit.Thickness,
		kit.Material, kit.Modifying, kit.Drawing, position.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *PositionKitRepo) Copy(ctx context.Context, targetPositionId string, position *position_api.CopyPosition) (string, error) {
	query := fmt.Sprintf(`INSERT INTO %s (id, position_id, type_id, type_code, construction_id, construction_code, count, size, thickness, 
		material, modifying, drawing) 
		SELECT $1, $2, type_id, type_code, construction_id, construction_code, count, size, thickness, 
		material, modifying, replace(drawing, $3, $4) FROM %s WHERE position_id=$5  RETURNING drawing`,
		PositionRingsKitTable, PositionRingsKitTable,
	)
	id := uuid.New()

	row := r.db.QueryRow(query, id, targetPositionId, position.FromOrderId, position.OrderId, position.Id)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", row.Err())
	}

	drawing := ""
	if err := row.Scan(&drawing); err != nil {
		return "", fmt.Errorf("failed to scan result design query. error: %w", row.Err())
	}

	return drawing, nil
}
