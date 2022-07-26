package pro_model

import "github.com/Alexander272/sealur_proto/api/pro_api"

type CreateAdditDTO struct {
	Materials    string `json:"materials" binding:"required"`
	Mod          string `json:"mod" binding:"required"`
	Temperature  string `json:"temperature" binding:"required"`
	Mounting     string `json:"mounting" binding:"required"`
	Graphite     string `json:"graphite" binding:"required"`
	Fillers      string `json:"fillers" binding:"required"`
	Coating      string `json:"coating" binding:"required"`
	Construction string `json:"construction" binding:"required"`
	Obturator    string `json:"obturator" binding:"required"`
	Basis        string `json:"basis" binding:"required"`
	PObt         string `json:"pobturator" binding:"required"`
	Sealant      string `json:"sealant" binding:"required"`
}

type UpdateMatDTO struct {
	Materials []*pro_api.AddMaterials `json:"materials" binding:"required"`
	TypeCh    string                  `json:"type" binding:"required"`
	Change    string                  `json:"change"`
}

type UpdateModDTO struct {
	Mod    []*pro_api.AddMod `json:"mod" binding:"required"`
	TypeCh string            `json:"type" binding:"required"`
	Change string            `json:"change"`
}

type UpdateTempDTO struct {
	Temperature []*pro_api.AddTemperature `json:"temperature" binding:"required"`
	TypeCh      string                    `json:"type" binding:"required"`
	Change      string                    `json:"change"`
}

type UpdateMounDTO struct {
	Mounting []*pro_api.AddMoun `json:"mounting" binding:"required"`
	TypeCh   string             `json:"type" binding:"required"`
	Change   string             `json:"change"`
}

type UpdateGrapDTO struct {
	Graphite []*pro_api.AddGrap `json:"graphite" binding:"required"`
	TypeCh   string             `json:"type" binding:"required"`
	Change   string             `json:"change"`
}

type UpdateFillersDTO struct {
	Fillers []*pro_api.AddFiller `json:"fillers" binding:"required"`
	TypeCh  string               `json:"type" binding:"required"`
	Change  string               `json:"change"`
}

type UpdateCoatingDTO struct {
	Coating []*pro_api.AddCoating `json:"coating" binding:"required"`
	TypeCh  string                `json:"type" binding:"required"`
	Change  string                `json:"change"`
}

type UpdateConstrDTO struct {
	Constr []*pro_api.AddConstruction `json:"constr" binding:"required"`
	TypeCh string                     `json:"type" binding:"required"`
	Change string                     `json:"change"`
}

type UpdateObturatorDTO struct {
	Obt    []*pro_api.AddObturator `json:"obturator" binding:"required"`
	TypeCh string                  `json:"type" binding:"required"`
	Change string                  `json:"change"`
}

type UpdateBasisDTO struct {
	Basis  []*pro_api.AddBasis `json:"basis" binding:"required"`
	TypeCh string              `json:"type" binding:"required"`
	Change string              `json:"change"`
}

type UpdatePObtDTO struct {
	Obturator []*pro_api.AddPObturator `json:"pobturator" binding:"required"`
	TypeCh    string                   `json:"type" binding:"required"`
	Change    string                   `json:"change"`
}

type UpdateSealantDTO struct {
	Sealant []*pro_api.AddSealant `json:"sealant" binding:"required"`
	TypeCh  string                `json:"type" binding:"required"`
	Change  string                `json:"change"`
}
