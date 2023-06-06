package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/flange_standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_configuration_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_construction_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_filler_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_flange_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_material_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_size_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PositionPutgRepo struct {
	db *sqlx.DB
}

func NewPositionPutgRepo(db *sqlx.DB) *PositionPutgRepo {
	return &PositionPutgRepo{
		db: db,
	}
}

func (r *PositionPutgRepo) Get(ctx context.Context, orderId string) (positions []*position_model.FullPosition, err error) {
	var data []models.PutgPosition
	query := fmt.Sprintf(`SELECT p.id, title, amount, type, count, info,
		configuration_id, configuration_code,
		filler_code, type_code, construction_code, reinforce_code, rotary_plug_code, inner_ring_code, outer_ring_code,
		d4, d3, d2, d1, h, use_dimensions,
		has_jumper, jumper_code, jumper_width, has_hole, has_coating, has_removable, has_mounting, mounting_code, drawing
		FROM %s AS p 
		INNER JOIN %s AS pm ON p.id=pm.position_id
		INNER JOIN %s AS pmt ON p.id=pmt.position_id
		INNER JOIN %s AS ps ON p.id=pmt.position_id
		INNER JOIN %s AS pd ON p.id=pd.position_id
		WHERE order_id=$1 AND type=$2 ORDER BY count`,
		PositionTable, PositionMainPutgTable, PositionMaterialPutgTable, PositionSizePutgTable, PositionDesignPutgTable,
	)

	if err := r.db.Select(&data, query, orderId, position_model.PositionType_Putg.String()); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, pp := range data {
		positions = append(positions, &position_model.FullPosition{
			Id:     pp.Id,
			Count:  pp.Count,
			Title:  pp.Title,
			Amount: pp.Amount,
			Info:   pp.Info,
			Type:   position_model.PositionType_Putg,
			PutgData: &position_model.PositionPutg{
				Main: &position_model.PositionPutg_Main{
					ConfigurationId:   pp.ConfigurationId,
					ConfigurationCode: pp.ConfigurationCode,
				},
				Material: &position_model.PositionPutg_Material{
					FillerCode:       pp.FillerCode,
					TypeCode:         pp.TypeCode,
					ConstructionCode: pp.ConstructionCode,
					ReinforceCode:    pp.ReinforceCode,
					RotaryPlugCode:   pp.RotaryPlugCode,
					InnerRindCode:    pp.InnerRingCode,
					OuterRingCode:    pp.OuterRingCode,
				},
				Size: &position_model.PositionPutg_Size{
					D4:            pp.D4,
					D3:            pp.D3,
					D2:            pp.D2,
					D1:            pp.D1,
					H:             pp.H,
					UseDimensions: pp.UseDimensions,
				},
				Design: &position_model.PositionPutg_Design{
					HasJumper:    pp.HasJumper,
					JumperCode:   pp.JumperCode,
					JumperWidth:  pp.JumperWidth,
					HasHole:      pp.HasHole,
					HasCoating:   pp.HasCoating,
					HasRemovable: pp.HasRemovable,
					HasMounting:  pp.HasMounting,
					MountingCode: pp.MountingCode,
					Drawing:      pp.Drawing,
				},
			},
		})
	}

	return positions, nil
}

func (r *PositionPutgRepo) GetFull(ctx context.Context, positionsId []string) ([]*position_model.OrderPositionPutg, error) {
	var mainData []models.PutgMainBlock
	mainQuery := fmt.Sprintf(`SELECT main.id, position_id, 
		main.putg_standard_id, ps.dn_title, ps.pn_title,
		ps.standard_id, s.title as standard_title,
		flange_standard_id, fs.title as flange_title, fs.code as flange_code,
		flange_type_id, ft.code as flange_type_code, ft.title as flange_type_title,
		configuration_id, conf.title as conf_title, conf.code as conf_code, conf.has_standard as conf_has_standard, conf.has_drawing as conf_has_drawing
		FROM %s as main
		LEFT JOIN %s AS ps ON ps.id=main.putg_standard_id
		LEFT JOIN %s AS s ON s.id=ps.standard_id
		LEFT JOIN %s AS fs ON fs.id=ps.flange_standard_id
		LEFT JOIN %s AS ft ON ft.id=flange_type_id
		LEFT JOIN %s AS conf ON conf.id=configuration_id
		WHERE array[position_id] <@ $1 ORDER BY position_id`,
		PositionMainPutgTable, PutgStandardTable, StandardTable, FlangeStandardTable, PutgFlangeTypeTable, PutgConfTable,
	)
	if err := r.db.Select(&mainData, mainQuery, pq.Array(positionsId)); err != nil {
		return nil, fmt.Errorf("failed to complete query main. error: %w", err)
	}

	var materialData []models.PutgMaterialBlock
	materialQuery := fmt.Sprintf(`SELECT pm.id, position_id, 
		pm.filler_id, pf.base_filler_id as f_base_id, f.code as f_code, f.title as f_title, f.description as f_description, f.designation as f_designation,
		type_id, t.title as t_title, t.code as t_code, t.description as t_description, t.min_thickness as t_min, t.max_thickness as t_max, 
		t.type_code as t_type_code, t.has_reinforce as t_has_reinforce,
		pm.construction_id, construction_code, pc.construction_id as c_base_id, c.title as c_title, c.description as c_description, c.has_d4 as c_has_d4, 
		c.has_d3 as c_has_d3, c.has_d2 as c_has_d2, c.has_d1 as c_has_d1, c.has_rotary_plug as c_has_rotary_plug, 
		c.has_inner_ring as c_has_inner_ring, c.has_outer_ring as c_has_outer_ring,
		reinforce_id, reinforce_code, reinforce_title, r.code as r_code, r.material_id as r_material_id, r.type as r_type,
		r.is_default as r_is_default,
		rotary_plug_id, rotary_plug_code, rotary_plug_title, rp.code as rp_code, rp.material_id as rp_material_id, rp.type as rp_type, 
		rp.is_default as rp_is_default,  
		inner_ring_id, inner_ring_code, inner_ring_title, ir.code as ir_code, ir.material_id as ir_material_id, ir.type as ir_type, 
		ir.is_default as ir_is_default,
		outer_ring_id, outer_ring_code, outer_ring_title, our.code as or_code, our.material_id as or_material_id, our.type as or_type, 
		our.is_default as or_is_default
		FROM %s AS pm
		LEFT JOIN %s AS pf ON filler_id=pf.id
		LEFT JOIN %s AS f ON base_filler_id=f.id
		LEFT JOIN %s AS t ON type_id=t.id
		LEFT JOIN %s AS pc ON pm.construction_id=pc.id
		LEFT JOIN %s AS c ON pc.construction_id=c.id 
		LEFT JOIN %s as r ON r.id=reinforce_id
		LEFT JOIN %s as rp ON rp.id=rotary_plug_id
		LEFT JOIN %s as ir ON ir.id=inner_ring_id
		LEFT JOIN %s as our ON our.id=outer_ring_id
		WHERE array[position_id] <@ $1 ORDER BY position_id`,
		PositionMaterialPutgTable, PutgFillerTable, PutgFillerBaseTable, PutgTypeTable, PutgConstructionTable, PutgConstructionBaseTable,
		PutgMaterialTable, PutgMaterialTable, PutgMaterialTable, PutgMaterialTable,
	)
	if err := r.db.Select(&materialData, materialQuery, pq.Array(positionsId)); err != nil {
		return nil, fmt.Errorf("failed to complete material query. error: %w", err)
	}

	var sizeData []models.PutgSizeBlock
	sizeQuery := fmt.Sprintf(`SELECT id, position_id, dn, dn_mm, pn_mpa, pn_kg, d4, d3, d2, d1, h, another, use_dimensions
		FROM %s WHERE array[position_id] <@ $1 ORDER BY position_id`,
		PositionSizePutgTable,
	)
	if err := r.db.Select(&sizeData, sizeQuery, pq.Array(positionsId)); err != nil {
		return nil, fmt.Errorf("failed to complete query size. error: %w", err)
	}

	var designData []models.PutgDesignBlock
	designQuery := fmt.Sprintf(`SELECT id, position_id, has_jumper, jumper_code, jumper_width, has_hole, has_coating, has_removable, 
		has_mounting, mounting_code, drawing FROM %s WHERE array[position_id] <@ $1 ORDER BY position_id`,
		PositionDesignPutgTable,
	)
	if err := r.db.Select(&designData, designQuery, pq.Array(positionsId)); err != nil {
		return nil, fmt.Errorf("failed to complete query size. error: %w", err)
	}

	positions := []*position_model.OrderPositionPutg{}
	for i, m := range mainData {
		mat := materialData[i]
		s := sizeData[i]
		d := designData[i]

		reinforce := &putg_material_model.Material{Id: mat.ReinforceId, BaseCode: mat.ReinforceBaseCode, Title: mat.ReinforceTitle}
		if mat.ReinforceId != uuid.Nil.String() {
			reinforce.Code = *mat.ReinforceCode
			reinforce.MaterialId = *mat.ReinforceMaterialId
			reinforce.Type = *mat.ReinforceType
			reinforce.IsDefault = *mat.ReinforceIsDefault
		} else {
			reinforce = nil
		}
		rotary := &putg_material_model.Material{Id: mat.RotaryPlugId, BaseCode: mat.RotaryPlugBaseCode, Title: mat.RotaryPlugTitle}
		if mat.RotaryPlugId != uuid.Nil.String() {
			rotary.Code = *mat.RotaryPlugCode
			rotary.MaterialId = *mat.RotaryPlugMaterialId
			rotary.Type = *mat.RotaryPlugType
			rotary.IsDefault = *mat.RotaryPlugIsDefault
		} else {
			rotary = nil
		}
		innerRing := &putg_material_model.Material{Id: mat.InnerRingId, BaseCode: mat.InnerRingBaseCode, Title: mat.InnerRingTitle}
		if mat.InnerRingId != uuid.Nil.String() {
			innerRing.Code = *mat.InnerRingCode
			innerRing.MaterialId = *mat.InnerRingMaterialId
			innerRing.Type = *mat.InnerRingType
			innerRing.IsDefault = *mat.InnerRingIsDefault
		} else {
			innerRing = nil
		}
		outerRing := &putg_material_model.Material{Id: mat.OuterRingId, BaseCode: mat.OuterRingBaseCode, Title: mat.OuterRingTitle}
		if mat.OuterRingId != uuid.Nil.String() {
			outerRing.Code = *mat.OuterRingCode
			outerRing.MaterialId = *mat.OuterRingMaterialId
			outerRing.Type = *mat.OuterRingType
			outerRing.IsDefault = *mat.OuterRingIsDefault
		} else {
			outerRing = nil
		}

		positions = append(positions, &position_model.OrderPositionPutg{
			Main: &position_model.OrderPositionPutg_Main{
				Id:         m.Id,
				PositionId: m.PositionId,
				Configuration: &putg_configuration_model.PutgConfiguration{
					Id:          m.ConfId,
					Title:       m.ConfTitle,
					Code:        m.ConfCode,
					HasStandard: m.ConfHasStandard,
					HasDrawing:  m.ConfHasDrawing,
				},
				FlangeType: &putg_flange_type_model.PutgFlangeType{
					Id:    m.FlangeTypeId,
					Title: m.FlangeTypeTitle,
					Code:  m.FlangeTypeCode,
				},
				Standard: &putg_standard_model.PutgStandard{
					Id:      m.PutgStandardId,
					DnTitle: m.PutgStandardDn,
					PnTitle: m.PutgStandardPn,
					Standard: &standard_model.Standard{
						Id:    m.StandardId,
						Title: m.StandardTitle,
					},
					FlangeStandard: &flange_standard_model.FlangeStandard{
						Id:    m.FlangeId,
						Code:  m.FlangeCode,
						Title: m.FlangeTitle,
					},
				},
			},
			Material: &position_model.OrderPositionPutg_Material{
				Id:         mat.Id,
				PositionId: mat.PositionId,
				Reinforce:  reinforce,
				RotaryPlug: rotary,
				InnerRing:  innerRing,
				OuterRing:  outerRing,
				Filler: &putg_filler_model.PutgFiller{
					Id:          mat.FillerId,
					BaseId:      mat.FillerBaseId,
					Title:       mat.FillerTitle,
					Code:        mat.FillerCode,
					Description: mat.FillerDescription,
					Designation: mat.FillerDesignation,
				},
				PutgType: &putg_type_model.PutgType{
					Id:           mat.TypeId,
					Title:        mat.TypeTitle,
					Code:         mat.TypeCode,
					TypeCode:     mat.TypeBaseCode,
					Description:  mat.TypeDescription,
					MinThickness: mat.TypeMinThickness,
					MaxThickness: mat.TypeMaxThickness,
					HasReinforce: mat.TypeHasReinforce,
				},
				Construction: &putg_construction_type_model.PutgConstruction{
					Id:            mat.ConstructionId,
					Title:         mat.ConstructionTitle,
					Code:          mat.ConstructionCode,
					Description:   mat.ConstructionDescription,
					BaseId:        mat.ConstructionBaseId,
					HasD4:         mat.ConstructionHasD4,
					HasD3:         mat.ConstructionHasD3,
					HasD2:         mat.ConstructionHasD2,
					HasD1:         mat.ConstructionHasD1,
					HasRotaryPlug: mat.ConstructionHasRotaryPlug,
					HasInnerRing:  mat.ConstructionHasInnerRing,
					HasOuterRing:  mat.ConstructionHasOuterRing,
				},
			},
			Size: &position_model.OrderPositionPutg_Size{
				Id:         s.Id,
				PositionId: s.PositionId,
				Dn:         s.Dn,
				DnMm:       s.DnMm,
				Pn: &snp_size_model.Pn{
					Mpa: s.PnMpa,
					Kg:  s.PnKg,
				},
				D4:            s.D4,
				D3:            s.D3,
				D2:            s.D2,
				D1:            s.D1,
				H:             s.H,
				Another:       s.Another,
				UseDimensions: s.UseDimensions,
			},
			Design: &position_model.OrderPositionPutg_Design{
				Id:         d.Id,
				PositionId: d.PositionId,
				Jumper: &position_model.OrderPositionPutg_Design_Jumper{
					HasJumper: d.HasJumper,
					Code:      d.JumperCode,
					Width:     d.JumperWidth,
				},
				HasHole:      d.HasHole,
				HasCoating:   d.HasCoating,
				HasRemovable: d.HasRemovable,
				Mounting: &position_model.OrderPositionPutg_Design_Mounting{
					HasMounting: d.HasMounting,
					Code:        d.MountingCode,
				},
				Drawing: d.Drawing,
			},
		})
	}

	return positions, nil
}

func (r *PositionPutgRepo) Create(ctx context.Context, position *position_model.FullPosition) error {
	mainQuery := fmt.Sprintf(`INSERT INTO %s (id, position_id, putg_standard_id, flange_type_id, configuration_id, configuration_code)
		VALUES ($1, $2, $3, $4, $5, $6)`, PositionMainPutgTable,
	)
	sizeQuery := fmt.Sprintf(`INSERT INTO %s (id, position_id, dn, dn_mm, pn_mpa, pn_kg, d4, d3, d2, d1, h, another, use_dimensions)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`, PositionSizePutgTable,
	)
	materialQuery := fmt.Sprintf(`INSERT INTO %s (id, position_id, filler_id, filler_code, type_id, type_code, construction_id, 
		construction_code, rotary_plug_id, rotary_plug_code, inner_ring_id, inner_ring_code, outer_ring_id, outer_ring_code, 
		rotary_plug_title, inner_ring_title, outer_ring_title, reinforce_id, reinforce_code, reinforce_title)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`,
		PositionMaterialPutgTable,
	)
	designQuery := fmt.Sprintf(`INSERT INTO %s (id, position_id, has_jumper, jumper_code, jumper_width, has_hole, has_coating, has_removable, 
		has_mounting, mounting_code, drawing) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, PositionDesignPutgTable,
	)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction. error: %w", err)
	}

	id := uuid.New()
	main := position.PutgData.Main
	size := position.PutgData.Size
	material := position.PutgData.Material
	design := position.PutgData.Design

	if material.ReinforceId == "" {
		material.ReinforceId = uuid.Nil.String()
	}
	if material.RotaryPlugId == "" {
		material.RotaryPlugId = uuid.Nil.String()
	}
	if material.InnerRingId == "" {
		material.InnerRingId = uuid.Nil.String()
	}
	if material.OuterRingId == "" {
		material.OuterRingId = uuid.Nil.String()
	}

	_, err = tx.Exec(mainQuery, id, position.Id, main.PutgStandardId, main.FlangeTypeId, main.ConfigurationId, main.ConfigurationCode)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query main. error: %w", err)
	}

	_, err = tx.Exec(sizeQuery, id, position.Id, size.Dn, size.DnMm, size.Pn.Mpa, size.Pn.Kg, size.D4, size.D3, size.D2, size.D1,
		size.H, size.Another, size.UseDimensions,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query size. error: %w", err)
	}

	_, err = tx.Exec(materialQuery, id, position.Id, material.FillerId, material.FillerCode, material.TypeId, material.TypeCode,
		material.ConstructionId, material.ConstructionCode, material.RotaryPlugId, material.RotaryPlugCode,
		material.InnerRingId, material.InnerRindCode, material.OuterRingId, material.OuterRingCode,
		material.RotaryPlugTitle, material.InnerRingTitle, material.OuterRingTitle, material.ReinforceId, material.ReinforceCode, material.ReinforceTitle,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query material. error: %w", err)
	}

	_, err = tx.Exec(designQuery, id, position.Id, design.HasJumper, design.JumperCode, design.JumperWidth, design.HasHole, design.HasCoating,
		design.HasRemovable, design.HasMounting, design.MountingCode, design.Drawing,
	)
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

//? нужен ли тут CreateSeveral

func (r *PositionPutgRepo) Update(ctx context.Context, position *position_model.FullPosition) error {
	mainQuery := fmt.Sprintf(`UPDATE %s SET putg_standard_id=$1, flange_type_id=$2, configuration_id=$3, configuration_code=$4 WHERE position_id=$5`,
		PositionMainPutgTable,
	)
	sizeQuery := fmt.Sprintf(`UPDATE %s SET dn=$1, dn_mm=$2, pn_mpa=$3, pn_kg=$4, d4=$5, d3=$6, d2=$7, d1=$8, h=$9, another=$10, use_dimensions=$11
		WHERE position_id=$12`,
		PositionSizePutgTable,
	)
	materialQuery := fmt.Sprintf(`UPDATE %s SET filler_id=$1, filler_code=$2, type_id=$3, type_code=$4, construction_id=$5, construction_code=$6, 
		rotary_plug_id=$7, rotary_plug_code=$8, inner_ring_id=$9, inner_ring_code=$10, outer_ring_id=$11, outer_ring_code=$12,
		rotary_plug_title=$13, inner_ring_title=$14, outer_ring_title=$15, reinforce_id=$16, reinforce_code=$17, reinforce_title=$18
		WHERE position_id=$19`,
		PositionMaterialPutgTable,
	)
	designQuery := fmt.Sprintf(`UPDATE %s SET has_jumper=$1, jumper_code=$2, jumper_width=$3, has_hole=$4, has_coating=$5, has_removable=$6, 
		has_mounting=$7, mounting_code=$8, drawing=$9 WHERE position_id=$10`,
		PositionDesignPutgTable,
	)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction. error: %w", err)
	}

	main := position.PutgData.Main
	size := position.PutgData.Size
	material := position.PutgData.Material
	design := position.PutgData.Design

	if material.ReinforceId == "" {
		material.ReinforceId = uuid.Nil.String()
	}
	if material.RotaryPlugId == "" {
		material.RotaryPlugId = uuid.Nil.String()
	}
	if material.InnerRingId == "" {
		material.InnerRingId = uuid.Nil.String()
	}
	if material.OuterRingId == "" {
		material.OuterRingId = uuid.Nil.String()
	}

	_, err = tx.Exec(mainQuery, main.PutgStandardId, main.FlangeTypeId, main.ConfigurationId, main.ConfigurationCode, position.Id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete main query. error: %w", err)
	}

	_, err = tx.Exec(sizeQuery, size.Dn, size.DnMm, size.Pn.Mpa, size.Pn.Kg, size.D4, size.D3, size.D2, size.D1, size.H,
		size.Another, size.UseDimensions, position.Id,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete size query. error: %w", err)
	}

	_, err = tx.Exec(materialQuery, material.FillerId, material.FillerCode, material.TypeId, material.TypeCode,
		material.ConstructionId, material.ConstructionCode,
		material.RotaryPlugId, material.RotaryPlugCode, material.InnerRingId, material.InnerRindCode, material.OuterRingId, material.OuterRingCode,
		material.RotaryPlugTitle, material.InnerRingTitle, material.OuterRingTitle, material.ReinforceId, material.ReinforceCode,
		material.ReinforceTitle, position.Id,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete material query. error: %w", err)
	}

	_, err = tx.Exec(designQuery, design.HasJumper, design.JumperCode, design.JumperWidth, design.HasHole, design.HasCoating,
		design.HasRemovable, design.HasMounting, design.MountingCode, design.Drawing, position.Id,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete design query. error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to finish transaction. error: %w", err)
	}
	return nil
}

func (r *PositionPutgRepo) Copy(ctx context.Context, targetPositionId string, position *position_api.CopyPosition) (string, error) {
	mainQuery := fmt.Sprintf(`INSERT INTO %s (id, position_id, putg_standard_id, flange_type_id, configuration_id, configuration_code)
		SELECT $1, $2, putg_standard_id, flange_type_id, configuration_id, configuration_code FROM %s WHERE position_id=$3`,
		PositionMainPutgTable, PositionMainPutgTable,
	)
	sizeQuery := fmt.Sprintf(`INSERT INTO %s (id, position_id, dn, dn_mm, pn_mpa, pn_kg, d4, d3, d2, d1, h, another, use_dimensions)
		SELECT $1, $2, dn, dn_mm, pn_mpa, pn_kg, d4, d3, d2, d1, h, another, use_dimensions FROM %s WHERE position_id=$3`,
		PositionSizePutgTable, PositionSizePutgTable,
	)
	materialQuery := fmt.Sprintf(`INSERT INTO %s (id, position_id, filler_id, filler_code, type_id, type_code, construction_id, construction_code,
		rotary_plug_id, rotary_plug_code, inner_ring_id, inner_ring_code, outer_ring_id, outer_ring_code, rotary_plug_title,
		inner_ring_title, outer_ring_title, reinforce_id, reinforce_code, reinforce_title)
		SELECT $1, $2, filler_id, filler_code, type_id, type_code, construction_id, construction_code,
		rotary_plug_id, rotary_plug_code, inner_ring_id, inner_ring_code, outer_ring_id, outer_ring_code, rotary_plug_title,
		inner_ring_title, outer_ring_title, reinforce_id, reinforce_code, reinforce_title FROM %s WHERE position_id=$3`,
		PositionMaterialPutgTable, PositionMaterialPutgTable,
	)
	designQuery := fmt.Sprintf(`INSERT INTO %s (id, position_id, has_jumper, jumper_code, jumper_width, has_hole, has_coating, has_removable,
		has_mounting, mounting_code, drawing) 
		SELECT $1, $2, has_jumper, jumper_code, jumper_width, has_hole, has_coating, has_removable, has_mounting, mounting_code, replace(drawing, $3, $4)  
		FROM %s WHERE position_id=$5 RETURNING drawing`,
		PositionDesignPutgTable, PositionDesignPutgTable,
	)

	id := uuid.New()

	tx, err := r.db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to start transaction. error: %w", err)
	}

	_, err = tx.Exec(mainQuery, id, targetPositionId, position.Id)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to complete main query. error: %w", err)
	}

	_, err = tx.Exec(materialQuery, id, targetPositionId, position.Id)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to complete material query. error: %w", err)
	}

	_, err = tx.Exec(sizeQuery, id, targetPositionId, position.Id)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to complete size query. error: %w", err)
	}

	row := tx.QueryRow(designQuery, id, targetPositionId, position.FromOrderId, position.OrderId, position.Id)
	if row.Err() != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to complete design query. error: %w", row.Err())
	}

	drawing := ""
	if err := row.Scan(&drawing); err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to scan result design query. error: %w", row.Err())
	}

	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("failed to finish transaction. error: %w", err)
	}
	return drawing, nil
}

func (r *PositionPutgRepo) Delete(ctx context.Context, positionId string) error {
	mainQuery := fmt.Sprintf(`DELETE FROM %s WHERE position_id=$1`, PositionMainPutgTable)
	sizeQuery := fmt.Sprintf(`DELETE FROM %s WHERE position_id=$1`, PositionSizePutgTable)
	materialQuery := fmt.Sprintf(`DELETE FROM %s WHERE position_id=$1`, PositionMaterialPutgTable)
	designQuery := fmt.Sprintf(`DELETE FROM %s WHERE position_id=$1`, PositionDesignPutgTable)

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
