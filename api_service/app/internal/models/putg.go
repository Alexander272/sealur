package models

import "github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"

type PutgDTO struct {
	FlangeId     string                     `json:"flangeId" binding:"required"`
	TypeFlId     string                     `json:"typeFlId" binding:"required"`
	TypePr       string                     `json:"typePr" binding:"required"`
	Form         string                     `json:"form" binding:"required"`
	Construction []*proto.PutgConstructions `json:"construction" binding:"required"`
	Temperatures []*proto.PutgTemp          `json:"temperatures" binding:"required"`
	Reinforce    *proto.PutgMaterials       `json:"reinforce"`
	Obturator    *proto.PutgMaterials       `json:"obturator"`
	ILimiter     *proto.PutgMaterials       `json:"iLimiter"`
	OLimiter     *proto.PutgMaterials       `json:"oLimiter"`
	Coating      []string                   `json:"coating" binding:"required"`
	Mounting     []string                   `json:"mounting" binding:"required"`
	Graphite     []string                   `json:"graphite" binding:"required"`
}
