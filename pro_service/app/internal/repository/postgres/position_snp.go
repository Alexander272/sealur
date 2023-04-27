package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/flange_standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_filler_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_material_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_size_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
	query := fmt.Sprintf(`SELECT %s.id, title, amount, info, type, count, filler_code, frame_code, inner_ring_code, outer_ring_code, d4, d3, d2, d1, h, another, 
		has_jumper, jumper_code, jumper_width, has_hole, has_mounting, mounting_code, drawing
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
			Info:   fp.Info,
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
					Drawing:      fp.Drawing,
				},
			},
		})
	}

	return positions, nil
}

// ? возможно стоит попробовать написать запросы получше
// получение всех данных о позициях снп
func (r *PositionSnpRepo) GetFull(ctx context.Context, positionsId []string) ([]*position_model.OrderPositionSnp, error) {
	var materialData []models.SnpMaterialBlockNew
	materialQuery := fmt.Sprintf(`SELECT %s.id, position_id,
			filler_id, %s.code, %s.title, %s.base_code, %s.description, %s.designation, %s.disabled_types,
			frame_code, frame_id, frame_title, m1.code as m1_code, m1.material_id as m1_material_id, m1.type as m1_type, m1.is_default as m1_is_default, m1.is_standard as m1_is_standard,
			inner_ring_code, inner_ring_id, inner_ring_title, m2.code as m2_code, m2.material_id as m2_material_id, m2.type as m2_type, m2.is_default as m2_is_default, m2.is_standard as m2_is_standard,
			outer_ring_code, outer_ring_id, outer_ring_title, m3.code as m3_code, m3.material_id as m3_material_id, m3.type as m3_type, m3.is_default as m3_is_default, m3.is_standard as m3_is_standard
			FROM %s
			LEFT JOIN %s ON %s.id=filler_id
			LEFT JOIN %s as m1 ON m1.id=frame_id
			LEFT JOIN %s as m2 ON m2.id=inner_ring_id
			LEFT JOIN %s as m3 ON m3.id=outer_ring_id
			WHERE array[position_id] <@ $1 ORDER BY position_id`,
		PositionMaterialSnpTable,
		SnpFillerNewTable, SnpFillerNewTable, SnpFillerNewTable, SnpFillerNewTable, SnpFillerNewTable, SnpFillerNewTable,
		PositionMaterialSnpTable, SnpFillerNewTable, SnpFillerNewTable, SnpMaterialTableNew, SnpMaterialTableNew, SnpMaterialTableNew,
	)
	if err := r.db.Select(&materialData, materialQuery, pq.Array(positionsId)); err != nil {
		return nil, fmt.Errorf("failed to complete query material. error: %w", err)
	}

	var mainData []models.SnpMainBlock
	mainQuery := fmt.Sprintf(`SELECT %s.id, position_id, flange_type_code, flange_type_title, 
		snp_type_id, st.code, st.title, st.has_d4, st.has_d3, st.has_d2, st.has_d1,
		ss.id as snp_standard_id, ss.dn_title, ss.pn_title, ss.has_d2 as st_has_d2,
		standard_id, s.title as standard_title, s.format as standard_format, 
		flange_standard_id, fs.title as flange_title, fs.code as flange_code
		FROM %s LEFT JOIN %s as st ON st.id=snp_type_id
		LEFT JOIN %s as ss ON ss.id=%s.snp_standard_id
		LEFT JOIN %s as s ON s.id=standard_id LEFT JOIN %s as fs ON fs.id=flange_standard_id
		WHERE array[position_id] <@ $1 ORDER BY position_id`,
		PositionMainSnpTable, PositionMainSnpTable, SnpTypeTable, SnpStandardTable, PositionMainSnpTable, StandardTable, FlangeStandardTable,
	)
	if err := r.db.Select(&mainData, mainQuery, pq.Array(positionsId)); err != nil {
		return nil, fmt.Errorf("failed to complete query main. error: %w", err)
	}

	var sizeData []models.SnpSizeBlock
	sizeQuery := fmt.Sprintf(`SELECT id, position_id, dn, dn_mm, pn_mpa, pn_kg, d4, d3, d2, d1, h, s2, s3, another
		FROM %s WHERE array[position_id] <@ $1 ORDER BY position_id`, PositionSizeSnpTable)
	if err := r.db.Select(&sizeData, sizeQuery, pq.Array(positionsId)); err != nil {
		return nil, fmt.Errorf("failed to complete query size. error: %w", err)
	}

	var designData []models.SnpDesignBlock
	designQuery := fmt.Sprintf(`SELECT id, position_id, has_jumper, jumper_code, jumper_width, has_hole, has_mounting, mounting_code, drawing
		FROM %s WHERE array[position_id] <@ $1 ORDER BY position_id`, PositionDesignSnpTable)
	if err := r.db.Select(&designData, designQuery, pq.Array(positionsId)); err != nil {
		return nil, fmt.Errorf("failed to complete query size. error: %w", err)
	}

	positions := []*position_model.OrderPositionSnp{}
	for i, m := range mainData {
		material := materialData[i]
		size := sizeData[i]
		design := designData[i]

		frame := &snp_material_model.Material{Id: material.FrameId, BaseCode: material.FrameBaseCode, Title: material.FrameTitle}
		if material.FrameId != uuid.Nil.String() {
			frame.Code = *material.FrameCode
			frame.MaterialId = *material.FrameMaterialId
			frame.Type = *material.FrameType
			frame.IsDefault = *material.FrameIsDefault
			frame.IsStandard = *material.FrameIsStandard
		} else {
			frame = nil
		}
		innerRing := &snp_material_model.Material{Id: material.InnerRingId, BaseCode: material.InnerRingBaseCode, Title: material.InnerRingTitle}
		if material.InnerRingId != uuid.Nil.String() {
			innerRing.Code = *material.InnerRingCode
			innerRing.MaterialId = *material.InnerRingMaterialId
			innerRing.Type = *material.InnerRingType
			innerRing.IsDefault = *material.InnerRingIsDefault
			innerRing.IsStandard = *material.InnerRingIsStandard
		} else {
			innerRing = nil
		}
		outerRing := &snp_material_model.Material{Id: material.OuterRingId, BaseCode: material.OuterRingBaseCode, Title: material.OuterRingTitle}
		if material.OuterRingId != uuid.Nil.String() {
			outerRing.Code = *material.OuterRingCode
			outerRing.MaterialId = *material.OuterRingMaterialId
			outerRing.Type = *material.OuterRingType
			outerRing.IsDefault = *material.OuterRingIsDefault
			outerRing.IsStandard = *material.OuterRingIsStandard
		} else {
			outerRing = nil
		}

		positions = append(positions, &position_model.OrderPositionSnp{
			Main: &position_model.OrderPositionSnp_Main{
				Id:            m.Id,
				PositionId:    m.PositionId,
				SnpStandardId: m.SnpStandardId,
				SnpTypeId:     m.SnpTypeId,
				SnpType: &snp_type_model.SnpType{
					Id:    m.SnpTypeId,
					Title: m.SnpTypeTitle,
					Code:  m.SnpTypeCode,
					HasD4: m.SnpTypeHasD4,
					HasD3: m.SnpTypeHasD3,
					HasD2: m.SnpTypeHasD2,
					HasD1: m.SnpTypeHasD1,
				},
				FlangeTypeCode:  m.FlangeTypeCode,
				FlangeTypeTitle: m.FlangeTypeTitle,
				SnpStandard: &snp_standard_model.SnpStandard{
					Id:      m.SnpStandardId,
					DnTitle: m.SnpStandardDn,
					PnTitle: m.SnpStandardPn,
					HasD2:   m.SnpStandardHasD2,
					Standard: &standard_model.Standard{
						Id:     m.StandardId,
						Title:  m.StandardTitle,
						Format: m.StandardFormat,
					},
					FlangeStandard: &flange_standard_model.FlangeStandard{
						Id:    m.FlangeId,
						Title: m.FlangeTitle,
						Code:  m.FlangeCode,
					},
				},
			},
			Material: &position_model.OrderPositionSnp_Material{
				Id:         material.Id,
				PositionId: material.PositionId,
				Filler: &snp_filler_model.SnpFiller{
					Id:            material.FillerId,
					BaseCode:      material.FillerBaseCode,
					Code:          material.FillerCode,
					Title:         material.FillerTitle,
					Description:   material.FillerDescription,
					Designation:   material.FillerDesignation,
					DisabledTypes: material.FillerDisabledTypes,
				},
				Frame:     frame,
				InnerRing: innerRing,
				OuterRing: outerRing,
			},
			Size: &position_model.OrderPositionSnp_Size{
				Id:         size.Id,
				PositionId: size.PositionId,
				Dn:         size.Dn,
				DnMm:       size.DnMm,
				Pn: &snp_size_model.Pn{
					Mpa: size.PnMpa,
					Kg:  size.PnKg,
				},
				D4:      size.D4,
				D3:      size.D3,
				D2:      size.D2,
				D1:      size.D1,
				H:       size.H,
				S2:      size.S2,
				S3:      size.S3,
				Another: size.Another,
			},
			Design: &position_model.OrderPositionSnp_Design{
				Id:         design.Id,
				PositionId: design.PositionId,
				Jumper: &position_model.OrderPositionSnp_Design_Jumper{
					HasJumper: design.HasJumper,
					Code:      design.JumperCode,
					Width:     design.JumperWidth,
				},
				HasHole: design.HasHole,
				Mounting: &position_model.OrderPositionSnp_Design_Mounting{
					HasMounting: design.HasMounting,
					Code:        design.MountingCode,
				},
				Drawing: design.Drawing,
			},
		})
	}

	return positions, nil
}

func (r *PositionSnpRepo) Create(ctx context.Context, position *position_model.FullPosition) error {
	mainQuery := fmt.Sprintf(`INSERT INTO %s(id, position_id, snp_standard_id, snp_type_id, flange_type_code, flange_type_title) 
		VALUES ($1, $2, $3, $4, $5, $6)`, PositionMainSnpTable)
	sizeQuery := fmt.Sprintf(`INSERT INTO %s(id, position_id, dn, dn_mm, pn_mpa, pn_kg, d4, d3, d2, d1, h, s2, s3, another) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`, PositionSizeSnpTable)
	materialQuery := fmt.Sprintf(`INSERT INTO %s(id, position_id, filler_id, frame_id, inner_ring_id, outer_ring_id, filler_code, frame_code, 
		inner_ring_code, outer_ring_code, frame_title, inner_ring_title, outer_ring_title) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		PositionMaterialSnpTable)
	designQuery := fmt.Sprintf(`INSERT INTO %s(id, position_id, has_jumper, jumper_code, jumper_width, has_hole, has_mounting, mounting_code, drawing) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, PositionDesignSnpTable)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction. error: %w", err)
	}

	id := uuid.New()
	main := position.SnpData.Main
	size := position.SnpData.Size
	material := position.SnpData.Material
	design := position.SnpData.Design

	nilId := uuid.Nil.String()
	if material.InnerRingId == "" {
		material.InnerRingId = nilId
	}
	if material.OuterRingId == "" {
		material.OuterRingId = nilId
	}

	_, err = tx.Exec(mainQuery, id, position.Id, main.SnpStandardId, main.SnpTypeId, main.FlangeTypeCode, main.FlangeTypeTitle)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query main. error: %w", err)
	}
	_, err = tx.Exec(sizeQuery, id, position.Id, size.Dn, size.DnMm, size.Pn.Mpa, size.Pn.Kg, size.D4, size.D3, size.D2, size.D1, size.H, size.S2, size.S3, size.Another)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query size. error: %w", err)
	}
	_, err = tx.Exec(materialQuery, id, position.Id, material.FillerId, material.FrameId, material.InnerRingId, material.OuterRingId, material.FillerCode,
		material.FrameCode, material.InnerRingCode, material.OuterRingCode, material.FrameTitle, material.InnerRingTitle, material.OuterRingTitle)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query material. error: %w", err)
	}
	_, err = tx.Exec(designQuery, id, position.Id, design.HasJumper, design.JumperCode, design.JumperWidth, design.HasHole, design.HasMounting,
		design.MountingCode, design.Drawing)
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

func (r *PositionSnpRepo) CreateSeveral(ctx context.Context, positions []*position_model.FullPosition) error {
	mainQuery := fmt.Sprintf(`INSERT INTO %s(id, position_id, snp_standard_id, snp_type_id, flange_type_code, flange_type_title) VALUES `, PositionMainSnpTable)
	sizeQuery := fmt.Sprintf(`INSERT INTO %s(id, position_id, dn, dn_mm, pn_mpa, pn_kg, d4, d3, d2, d1, h, s2, s3, another) VALUES `, PositionSizeSnpTable)
	materialQuery := fmt.Sprintf(`INSERT INTO %s(id, position_id, filler_id, frame_id, inner_ring_id, outer_ring_id, filler_code, frame_code, inner_ring_code, outer_ring_code, frame_title, inner_ring_title, outer_ring_title) VALUES `, PositionMaterialSnpTable)
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
	sizeCount := 14
	materialCount := 13
	designCount := 9

	nilId := uuid.Nil.String()

	for i, p := range positions {
		id := uuid.New()

		mainValues = append(mainValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			i*mainCount+1, i*mainCount+2, i*mainCount+3, i*mainCount+4, i*mainCount+5, i*mainCount+6),
		)
		main = p.SnpData.Main
		mainArgs = append(mainArgs, id, p.Id, main.SnpStandardId, main.SnpTypeId, main.FlangeTypeCode, main.FlangeTypeTitle)

		sizeValues = append(sizeValues, fmt.Sprintf(`($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)`,
			i*sizeCount+1, i*sizeCount+2, i*sizeCount+3, i*sizeCount+4, i*sizeCount+5, i*sizeCount+6, i*sizeCount+7, i*sizeCount+8,
			i*sizeCount+9, i*sizeCount+10, i*sizeCount+11, i*sizeCount+12, i*sizeCount+13, i*sizeCount+14,
		))
		size = p.SnpData.Size
		sizeArgs = append(sizeArgs, id, p.Id, size.Dn, size.DnMm, size.Pn.Mpa, size.Pn.Kg, size.D4, size.D3, size.D2, size.D1, size.H, size.S2, size.S3, size.Another)

		materialValues = append(materialValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*materialCount+1, i*materialCount+2, i*materialCount+3, i*materialCount+4, i*materialCount+5, i*materialCount+6, i*materialCount+7,
			i*materialCount+8, i*materialCount+9, i*materialCount+10, i*materialCount+11, i*materialCount+12, i*materialCount+13,
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
			material.FrameTitle, material.InnerRingTitle, material.OuterRingTitle,
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

func (r *PositionSnpRepo) Update(ctx context.Context, position *position_model.FullPosition) error {
	mainQuery := fmt.Sprintf(`UPDATE %s SET snp_standard_id=$1, snp_type_id=$2, flange_type_code=$3, flange_type_title=$4 WHERE position_id=$5`,
		PositionMainSnpTable)
	sizeQuery := fmt.Sprintf(`UPDATE %s	SET dn=$1, dn_mm=$2, pn_mpa=$3, pn_kg=$4, d4=$5, d3=$6, d2=$7, d1=$8, h=$9, s2=$10, s3=$11, another=$12
		WHERE position_id=$13`, PositionSizeSnpTable)
	materialQuery := fmt.Sprintf(`UPDATE %s SET filler_id=$1, frame_id=$2, inner_ring_id=$3, outer_ring_id=$4, filler_code=$5, 
		frame_code=$6, inner_ring_code=$7, outer_ring_code=$8, frame_title=$9, inner_ring_title=$10, outer_ring_title=$11 WHERE position_id=$12`, PositionMaterialSnpTable)
	designQuery := fmt.Sprintf(`UPDATE %s SET has_jumper=$1, jumper_code=$2, jumper_width=$3, has_hole=$4, has_mounting=$5, mounting_code=$6, drawing=$7
		WHERE position_id=$8`, PositionDesignSnpTable)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction. error: %w", err)
	}

	main := position.SnpData.Main
	size := position.SnpData.Size
	material := position.SnpData.Material
	design := position.SnpData.Design

	nilId := uuid.Nil.String()
	if material.InnerRingId == "" {
		material.InnerRingId = nilId
	}
	if material.OuterRingId == "" {
		material.OuterRingId = nilId
	}

	_, err = tx.Exec(mainQuery, main.SnpStandardId, main.SnpTypeId, main.FlangeTypeCode, main.FlangeTypeTitle, position.Id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query main. error: %w", err)
	}
	_, err = tx.Exec(sizeQuery, size.Dn, size.DnMm, size.Pn.Mpa, size.Pn.Kg, size.D4, size.D3, size.D2, size.D1, size.H, size.S2, size.S3, size.Another, position.Id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query size. error: %w", err)
	}
	_, err = tx.Exec(materialQuery, material.FillerId, material.FrameId, material.InnerRingId, material.OuterRingId, material.FillerCode, material.FrameCode,
		material.InnerRingCode, material.OuterRingCode, material.FrameTitle, material.InnerRingTitle, material.OuterRingTitle, position.Id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query material. error: %w", err)
	}
	_, err = tx.Exec(designQuery, design.HasJumper, design.JumperCode, design.JumperWidth, design.HasHole, design.HasMounting,
		design.MountingCode, design.Drawing, position.Id)
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

func (r *PositionSnpRepo) Copy(ctx context.Context, targetPositionId string, position *position_api.CopyPosition) (string, error) {
	mainQuery := fmt.Sprintf(`INSERT INTO %s (id, position_id, snp_standard_id, snp_type_id, flange_type_code, flange_type_title)
		SELECT $1, $2, snp_standard_id, snp_type_id, flange_type_code, flange_type_title FROM %s WHERE position_id=$3`, PositionMainSnpTable, PositionMainSnpTable,
	)
	materialQuery := fmt.Sprintf(`INSERT INTO %s (id, position_id, filler_id, frame_id, inner_ring_id, outer_ring_id, 
		filler_code, frame_code, inner_ring_code, outer_ring_code) SELECT $1, $2, filler_id, frame_id, inner_ring_id, outer_ring_id, 
		filler_code, frame_code, inner_ring_code, outer_ring_code FROM %s WHERE position_id=$3`,
		PositionMaterialSnpTable, PositionMaterialSnpTable,
	)
	sizeQuery := fmt.Sprintf(`INSERT INTO %s (id, position_id, dn, dn_mm, pn_mpa, pn_kg, d4, d3, d2, d1, h, s2, s3, another)
		SELECT $1, $2, dn, dn_mm, pn_mpa, pn_kg, d4, d3, d2, d1, h, s2, s3, another FROM %s WHERE position_id=$3`,
		PositionSizeSnpTable, PositionSizeSnpTable,
	)
	// надо что-то думать с чертежом (ели он есть)
	designQuery := fmt.Sprintf(`INSERT INTO %s(id, position_id, has_jumper, jumper_code, jumper_width, has_hole, has_mounting, mounting_code, drawing)
		SELECT $1, $2, has_jumper, jumper_code, jumper_width, has_hole, has_mounting, mounting_code, replace(drawing, $3, $4) 
		FROM %s WHERE position_id=$5 RETURNING drawing`,
		PositionDesignSnpTable, PositionDesignSnpTable,
	)

	id := uuid.New()

	tx, err := r.db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to start transaction. error: %w", err)
	}

	_, err = tx.Exec(mainQuery, id, targetPositionId, position.Id)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to complete query main. error: %w", err)
	}
	_, err = tx.Exec(materialQuery, id, targetPositionId, position.Id)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to complete query material. error: %w", err)
	}
	_, err = tx.Exec(sizeQuery, id, targetPositionId, position.Id)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to complete query size. error: %w", err)
	}
	row := tx.QueryRow(designQuery, id, targetPositionId, position.FromOrderId, position.OrderId, position.Id)
	if row.Err() != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to complete query design. error: %w", row.Err())
	}

	drawing := ""
	if err := row.Scan(&drawing); err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to scan result query design. error: %w", row.Err())
	}

	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("failed to finish transaction. error: %w", err)
	}
	return drawing, nil
}

func (r *PositionSnpRepo) Delete(ctx context.Context, positionId string) error {
	mainQuery := fmt.Sprintf(`DELETE FROM %s WHERE position_id=$1`, PositionMainSnpTable)
	sizeQuery := fmt.Sprintf(`DELETE FROM %s WHERE position_id=$1`, PositionSizeSnpTable)
	materialQuery := fmt.Sprintf(`DELETE FROM %s WHERE position_id=$1`, PositionMaterialSnpTable)
	designQuery := fmt.Sprintf(`DELETE FROM %s WHERE position_id=$1`, PositionDesignSnpTable)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction. error: %w", err)
	}

	_, err = tx.Exec(mainQuery, positionId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query main. error: %w", err)
	}
	_, err = tx.Exec(sizeQuery, positionId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query size. error: %w", err)
	}
	_, err = tx.Exec(materialQuery, positionId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query material. error: %w", err)
	}
	_, err = tx.Exec(designQuery, positionId)
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
