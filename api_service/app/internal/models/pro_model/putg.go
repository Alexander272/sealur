package pro_model

import "github.com/Alexander272/sealur_proto/api/pro_api"

type PutgDTO struct {
	FlangeId     string                       `json:"flangeId" binding:"required"`
	TypeFlId     string                       `json:"typeFlId" binding:"required"`
	TypePr       string                       `json:"typePr" binding:"required"`
	Form         string                       `json:"form" binding:"required"`
	Construction []*pro_api.PutgConstructions `json:"construction" binding:"required"`
	Temperatures []*pro_api.PutgTemp          `json:"temperatures" binding:"required"`
	Reinforce    *pro_api.PutgMaterials       `json:"reinforce"`
	Obturator    *pro_api.PutgMaterials       `json:"obturator"`
	ILimiter     *pro_api.PutgMaterials       `json:"iLimiter"`
	OLimiter     *pro_api.PutgMaterials       `json:"oLimiter"`
	Coating      []string                     `json:"coating" binding:"required"`
	Mounting     []string                     `json:"mounting" binding:"required"`
	Graphite     []string                     `json:"graphite" binding:"required"`
}
