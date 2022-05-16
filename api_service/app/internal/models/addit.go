package models

import "github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"

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
	Basis        string `json:"basis"`
	Sealant      string `json:"sealant"`
}

type UpdateMatDTO struct {
	Materials []*proto.AddMaterials `json:"materials" binding:"required"`
	TypeCh    string                `json:"type" binding:"required"`
	Change    string                `json:"change"`
}

type UpdateModDTO struct {
	Mod    []*proto.AddMod `json:"mod" binding:"required"`
	TypeCh string          `json:"type" binding:"required"`
	Change string          `json:"change"`
}

type UpdateTempDTO struct {
	Temperature []*proto.AddTemperature `json:"temperature" binding:"required"`
	TypeCh      string                  `json:"type" binding:"required"`
	Change      string                  `json:"change"`
}

type UpdateMounDTO struct {
	Mounting []*proto.AddMoun `json:"mounting" binding:"required"`
	TypeCh   string           `json:"type" binding:"required"`
	Change   string           `json:"change"`
}

type UpdateGrapDTO struct {
	Graphite []*proto.AddGrap `json:"graphite" binding:"required"`
	TypeCh   string           `json:"type" binding:"required"`
	Change   string           `json:"change"`
}

type UpdateFillersDTO struct {
	Fillers []*proto.AddFiller `json:"fillers" binding:"required"`
	TypeCh  string             `json:"type" binding:"required"`
	Change  string             `json:"change"`
}

type UpdateCoatingDTO struct {
	Coating []*proto.AddCoating `json:"coating" binding:"required"`
	TypeCh  string              `json:"type" binding:"required"`
	Change  string              `json:"change"`
}

type UpdateConstrDTO struct {
	Constr []*proto.AddConstruction `json:"constr" binding:"required"`
	TypeCh string                   `json:"type" binding:"required"`
	Change string                   `json:"change"`
}

type UpdateObturatorDTO struct {
	Obt    []*proto.AddObturator `json:"obturator" binding:"required"`
	TypeCh string                `json:"type" binding:"required"`
	Change string                `json:"change"`
}

type UpdateBasisDTO struct {
	Basis  []*proto.AddBasis `json:"basis" binding:"required"`
	TypeCh string            `json:"type" binding:"required"`
	Change string            `json:"change"`
}

type UpdateSealantDTO struct {
	Sealant []*proto.AddSealant `json:"sealant" binding:"required"`
	TypeCh  string              `json:"type" binding:"required"`
	Change  string              `json:"change"`
}
