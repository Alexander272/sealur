package pro_model

import (
	"github.com/Alexander272/sealur_proto/api/pro_api"
)

type SNPDTO struct {
	StandId  string                `json:"standId" binding:"required"`
	FlangeId string                `json:"flangeId" binding:"required"`
	TypeFlId string                `json:"typeFlId" binding:"required"`
	TypePr   string                `json:"typePr" binding:"required"`
	Fillers  []*pro_api.Filler     `json:"fillers" binding:"required"`
	Frame    *pro_api.SnpMaterials `json:"frame"`
	Ir       *pro_api.SnpMaterials `json:"ir"`
	Or       *pro_api.SnpMaterials `json:"or"`
	Mounting []string              `json:"mounting" binding:"required"`
	Graphite []string              `json:"graphite" binding:"required"`
}

type DefResponse struct {
	TypeFl []*pro_api.TypeFl     `json:"typeFl"`
	Snp    []*pro_api.SNP        `json:"snp"`
	Sizes  *pro_api.SizeResponse `json:"sizes"`
}
