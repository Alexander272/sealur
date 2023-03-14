package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PositionSnpRepo struct {
	db *sqlx.DB
}

func NewPositionSnpRepo(db *sqlx.DB) *PositionSnpRepo {
	return &PositionSnpRepo{db: db}
}

func (r *PositionSnpRepo) Get(ctx context.Context, orderId string) (positions []*position_model.FullPosition, err error) {
	//? можно не делать запрос в position, а в этом запросе забрать данные из всех 5 таблиц через inner join
	var data []models.FullPosition
	query := fmt.Sprintf(`SELECT %s.id, title, amount, type, count, filler_code, frame_code, inner_ring_code, outer_ring_code, d4, d3, d2, d1, h, another, 
		has_jumper, jumper_code, jumper_width, has_hole, has_mounting, mounting_code
		FROM %s	INNER JOIN %s ON %s.position_id=%s.id INNER JOIN %s ON %s.position_id=%s.id INNER JOIN %s ON %s.position_id=%s.id
		WHERE order_id=$1 ORDER BY count`,
		PositionTable, PositionTable, PositionMaterialSnpTable, PositionMaterialSnpTable, PositionTable,
		PositionSizeSnpTable, PositionSizeSnpTable, PositionTable, PositionDesignSnpTable, PositionDesignSnpTable, PositionTable,
	)

	if err := r.db.Select(&data, query, orderId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, fp := range data {
		positions = append(positions, &position_model.FullPosition{
			Count:  fp.Count,
			Title:  fp.Title,
			Amount: fp.Amount,
			Type:   position_model.PositionType_Snp,
			SnpData: &position_model.PositionSnp{
				Size: &position_model.PositionSnp_Size{
					D4:      fp.D4,
					D3:      fp.D3,
					D2:      fp.D2,
					D1:      fp.D1,
					H:       fp.H,
					Another: fp.Another,
				},
				Material: &position_model.PositionSnp_Material{
					FillerCode:    fp.FillerCode,
					FrameCode:     fp.FrameCode,
					InnerRingCode: fp.InnerRingCode,
					OuterRingCode: fp.OuterRingCode,
				},
				Design: &position_model.PositionSnp_Design{
					HasJumper:    fp.HasJumper,
					JumperCode:   fp.JumperCode,
					JumperWidth:  fp.JumperWidth,
					HasHole:      fp.HasHole,
					HasMounting:  fp.HasMounting,
					MountingCode: fp.MountingCode,
				},
			},
		})
	}

	return positions, nil
}

func (r *PositionSnpRepo) CreateSeveral(ctx context.Context, positions []*position_model.FullPosition) error {
	mainQuery := fmt.Sprintf(`INSERT INTO %s(id, position_id, snp_standard_id, snp_type_id, flange_type_code, flange_type_title) VALUES `, PositionMainSnpTable)
	sizeQuery := fmt.Sprintf(`INSERT INTO %s(id, position_id, dn, pn_mpa, pn_kg, d4, d3, d2, d1, h, s2, s3, another) VALUES `, PositionSizeSnpTable)
	materialQuery := fmt.Sprintf(`INSERT INTO %s(id, position_id, filler_id, frame_id, inner_ring_id, outer_ring_id, filler_code, frame_code, inner_ring_code, outer_ring_code) VALUES `, PositionMaterialSnpTable)
	designQuery := fmt.Sprintf(`INSERT INTO %s(id, position_id, has_jumper, jumper_code, jumper_width, has_hole, has_mounting, mounting_code, drawing) VALUES `, PositionDesignSnpTable)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction. error: %w", err)
	}

	mainArgs := make([]interface{}, 0, len(positions))
	mainValues := make([]string, 0, len(positions))
	sizeArgs := make([]interface{}, 0, len(positions))
	sizeValues := make([]string, 0, len(positions))
	materialArgs := make([]interface{}, 0, len(positions))
	materialValues := make([]string, 0, len(positions))
	designArgs := make([]interface{}, 0, len(positions))
	designValues := make([]string, 0, len(positions))

	var main *position_model.PositionSnp_Main
	var size *position_model.PositionSnp_Size
	var material *position_model.PositionSnp_Material
	var design *position_model.PositionSnp_Design

	mainCount := 6
	sizeCount := 13
	materialCount := 10
	designCount := 9

	nilId := uuid.Nil.String()

	for i, p := range positions {
		id := uuid.New()

		mainValues = append(mainValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			i*mainCount+1, i*mainCount+2, i*mainCount+3, i*mainCount+4, i*mainCount+5, i*mainCount+6),
		)
		main = p.SnpData.Main
		mainArgs = append(mainArgs, id, p.Id, main.SnpStandardId, main.SnpTypeId, main.FlangeTypeCode, main.FlangeTypeTitle)

		sizeValues = append(sizeValues, fmt.Sprintf(`($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)`,
			i*sizeCount+1, i*sizeCount+2, i*sizeCount+3, i*sizeCount+4, i*sizeCount+5, i*sizeCount+6, i*sizeCount+7, i*sizeCount+8,
			i*sizeCount+9, i*sizeCount+10, i*sizeCount+11, i*sizeCount+12, i*sizeCount+13,
		))
		size = p.SnpData.Size
		sizeArgs = append(sizeArgs, id, p.Id, size.Dn, size.Pn.Mpa, size.Pn.Kg, size.D4, size.D3, size.D2, size.D1, size.H, size.S2, size.S3, size.Another)

		materialValues = append(materialValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*materialCount+1, i*materialCount+2, i*materialCount+3, i*materialCount+4, i*materialCount+5, i*materialCount+6, i*materialCount+7,
			i*materialCount+8, i*materialCount+9, i*materialCount+10,
		))
		material = p.SnpData.Material
		if material.InnerRingId == "" {
			material.InnerRingId = nilId
		}
		if material.OuterRingId == "" {
			material.OuterRingId = nilId
		}
		materialArgs = append(materialArgs, id, p.Id, material.FillerId, material.FrameId, material.InnerRingId, material.OuterRingId,
			material.FillerCode, material.FrameCode, material.InnerRingCode, material.OuterRingCode,
		)

		designValues = append(designValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*designCount+1, i*designCount+2, i*designCount+3, i*designCount+4, i*designCount+5, i*designCount+6, i*designCount+7, i*designCount+8, i*designCount+9,
		))
		design = p.SnpData.Design
		designArgs = append(designArgs, id, p.Id, design.HasJumper, design.JumperCode, design.JumperWidth, design.HasHole, design.HasMounting, design.MountingCode, design.Drawing)
	}

	mainQuery += strings.Join(mainValues, ",")
	sizeQuery += strings.Join(sizeValues, ",")
	materialQuery += strings.Join(materialValues, ",")
	designQuery += strings.Join(designValues, ",")

	//TODO тут где-то ошибка: failed to create several position snp. error: failed to execute query. error: pq: неверный синтаксис для типа uuid: \"\"
	_, err = tx.Exec(mainQuery, mainArgs...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query main. error: %w", err)
	}
	_, err = tx.Exec(sizeQuery, sizeArgs...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query size. error: %w", err)
	}
	_, err = tx.Exec(materialQuery, materialArgs...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query material. error: %w", err)
	}
	_, err = tx.Exec(designQuery, designArgs...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query design. error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to finish transaction. error: %w", err)
	}
	return nil
}
