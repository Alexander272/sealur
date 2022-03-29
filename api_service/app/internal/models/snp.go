package models

import (
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
)

type SNPDTO struct {
	StandId  string `json:"standId" binding:"required"`
	FlangeId string `json:"flangeId" binding:"required"`
	TypeFlId string `json:"typeFlId" binding:"required"`
	TypePr   string `json:"typePr" binding:"required"`
	Fillers  string `json:"fillers" binding:"required"`
	Frame    string `json:"frame"`
	Ir       string `json:"ir"`
	Or       string `json:"or"`
	Mounting string `json:"mounting" binding:"required"`
	Graphite string `json:"graphite" binding:"required"`
}

type DefResponse struct {
	TypeFl []*proto.TypeFl     `json:"typeFl"`
	Snp    []*proto.SNP        `json:"snp"`
	Sizes  *proto.SizeResponse `json:"sizes"`
}
