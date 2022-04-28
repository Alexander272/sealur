package models

import "github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"

type PutgDTO struct {
	FlangeId     string                     `json:"flangeId" binding:"required"`
	TypeFlId     string                     `json:"typeFlId" binding:"required"`
	TypePr       string                     `json:"typePr" binding:"required"`
	Form         string                     `json:"form" binding:"required"`
	Construction []*proto.PutgConstructions `json:"construction" binding:"required"`
	Temperatures []*proto.PutgTemp          `json:"temperatures" binding:"required"`
	Reinforce    *proto.Materials           `json:"reinforce"`
	Obturator    *proto.Materials           `json:"obturator"`
	ILimiter     *proto.Materials           `json:"iLimiter"`
	OLimiter     *proto.Materials           `json:"oLimiter"`
	Coating      []string                   `json:"coating" binding:"required"`
	Mounting     []string                   `json:"mounting" binding:"required"`
	Graphite     []string                   `json:"graphite" binding:"required"`
}
