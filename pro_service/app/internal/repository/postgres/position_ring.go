package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/ring_construction_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/ring_density_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/ring_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PositionRingRepo struct {
	db *sqlx.DB
}

func NewPositionRingRepo(db *sqlx.DB) *PositionRingRepo {
	return &PositionRingRepo{
		db: db,
	}
}

type PositionRing interface {
	Get(ctx context.Context, orderId string) ([]*position_model.FullPosition, error)
	GetFull(ctx context.Context, positionsId []string) ([]*position_model.OrderPositionRing, error)
	Create(context.Context, *position_model.FullPosition) error
	Update(context.Context, *position_model.FullPosition) error
	Copy(ctx context.Context, targetPositionId string, position *position_api.CopyPosition) (string, error)
}

func (r *PositionRingRepo) Get(ctx context.Context, orderId string) (positions []*position_model.FullPosition, err error) {
	var data []models.RingPosition
	query := fmt.Sprintf(`SELECT p.id, title, amount, type, count, info,
		type_code, density_code, construction_code, construction_bc, size, thickness, material, modifying, drawing
		FROM %s AS p LEFT JOIN %s ON p.id=position_id
		WHERE order_id=$1 AND type=$2 ORDER BY count`,
		PositionTable, PositionRingTable,
	)

	if err := r.db.Select(&data, query, orderId, position_model.PositionType_Ring.String()); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, rp := range data {
		positions = append(positions, &position_model.FullPosition{
			Id:     rp.Id,
			Count:  rp.Count,
			Title:  rp.Title,
			Amount: rp.Amount,
			Info:   rp.Info,
			Type:   position_model.PositionType_Ring,
			RingData: &position_model.PositionRing{
				TypeCode:             rp.TypeCode,
				DensityCode:          rp.DensityCode,
				ConstructionCode:     rp.ConstructionCode,
				ConstructionBaseCode: rp.ConstructionBaseCode,
				Size:                 rp.Size,
				Thickness:            rp.Thickness,
				Material:             rp.Material,
				Modifying:            rp.Modifying,
				Drawing:              rp.Drawing,
			},
		})
	}

	return positions, nil
}

func (r *PositionRingRepo) GetFull(ctx context.Context, positionsId []string) ([]*position_model.OrderPositionRing, error) {
	var data []models.Ring
	query := fmt.Sprintf(`SELECT r.id, r.type_id, t.code as type_code, t.image as t_image, t.material_type as t_mt, t.has_rotary_plug as t_hrp, 
		t.has_density as t_hd, t.has_thickness as t_ht, t.designation as t_designation,
		density_id, COALESCE(d.code, '') as d_code, COALESCE(d.has_rotary_plug,false) as d_hrp,
		construction_code, construction_wrp, construction_bc,
		size, thickness, material, modifying, drawing, position_id
		FROM %s AS r 
		LEFT JOIN %s AS t ON type_id=t.id 
		LEFT JOIN %s AS d ON density_id=d.id::text
		WHERE array[position_id] <@ $1 ORDER BY position_id`,
		PositionRingTable, RingTypeTable, RingDensityTable,
	)

	if err := r.db.Select(&data, query, pq.Array(positionsId)); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	positions := []*position_model.OrderPositionRing{}
	for _, ring := range data {
		positions = append(positions, &position_model.OrderPositionRing{
			Id:         ring.Id,
			PositionId: ring.PositionId,
			RingType: &ring_type_model.RingType{
				Id:            ring.TypeId,
				Code:          ring.TypeCode,
				Image:         ring.TypeImage,
				MaterialType:  ring.TypeMT,
				HasRotaryPlug: ring.TypeHRP,
				HasDensity:    ring.TypeHD,
				HasThickness:  ring.TypeHT,
				Designation:   ring.TypeDesignation,
			},
			Density: &ring_density_model.RingDensity{
				Id:            ring.DensityId,
				Code:          ring.DensityCode,
				HasRotaryPlug: ring.DensityHRP,
			},
			Construction: &ring_construction_model.RingConstruction{
				Code:              ring.ConstructionCode,
				WithoutRotaryPlug: ring.ConstructionWRP,
				BaseCode:          ring.ConstructionBaseCode,
			},
			Size:      ring.Size,
			Thickness: ring.Thickness,
			Material:  ring.Material,
			Modifying: ring.Modifying,
			Drawing:   ring.Drawing,
		})
	}

	return positions, nil
}

func (r *PositionRingRepo) Create(ctx context.Context, position *position_model.FullPosition) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, position_id, type_id, type_code, density_id, density_code, construction_code, 
		construction_wrp, construction_bc, size, thickness, material, modifying, drawing) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		PositionRingTable,
	)
	id := uuid.New()

	ring := position.RingData
	_, err := r.db.Exec(query, id, position.Id, ring.TypeId, ring.TypeCode, ring.DensityId, ring.DensityCode, ring.ConstructionCode,
		ring.ConstructionWRP, ring.ConstructionBaseCode, ring.Size, ring.Thickness, ring.Material, ring.Modifying, ring.Drawing,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *PositionRingRepo) Update(ctx context.Context, position *position_model.FullPosition) error {
	query := fmt.Sprintf(`UPDATE %s	SET type_id=$1, type_code=$2, density_id=$3, density_code=$4, construction_code=$5, construction_wrp=$6,
		construction_bc=$7, size=$8, thickness=$9, material=$10, modifying=$11, drawing=$12 
		WHERE position_id=$13`,
		PositionRingTable,
	)

	ring := position.RingData
	_, err := r.db.Exec(query, ring.TypeId, ring.TypeCode, ring.DensityId, ring.DensityCode, ring.ConstructionCode, ring.ConstructionWRP,
		ring.ConstructionBaseCode, ring.Size, ring.Thickness, ring.Material, ring.Modifying, ring.Drawing, position.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *PositionRingRepo) Copy(ctx context.Context, targetPositionId string, position *position_api.CopyPosition) (string, error) {
	query := fmt.Sprintf(`INSERT INTO %s(id, position_id, type_id, type_code, density_id, density_code, construction_code, 
		construction_wrp, construction_bc, size, thickness, material, modifying, drawing) 
		SELECT $1, $2, type_id, type_code, density_id, density_code, construction_code, construction_wrp, construction_bc, size, thickness, 
		material, modifying, replace(drawing, $3, $4) FROM %s WHERE position_id=$5  RETURNING drawing`,
		PositionRingTable, PositionRingTable,
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
