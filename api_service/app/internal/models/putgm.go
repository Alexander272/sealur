package models

import "github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"

type PutgmDTO struct {
	FlangeId     string                      `json:"flangeId" binding:"required"`
	TypeFlId     string                      `json:"typeFlId" binding:"required"`
	TypePr       string                      `json:"typePr" binding:"required"`
	Form         string                      `json:"form" binding:"required"`
	Construction []*proto.PutgmConstructions `json:"construction" binding:"required"`
	Temperatures []*proto.PutgTemp           `json:"temperatures" binding:"required"`
	Basis        *proto.PutgMaterials        `json:"basis"`
	Obturator    *proto.PutgMaterials        `json:"obturator"`
	Coating      []string                    `json:"coating" binding:"required"`
	Mounting     []string                    `json:"mounting" binding:"required"`
	Graphite     []string                    `json:"graphite" binding:"required"`
}
