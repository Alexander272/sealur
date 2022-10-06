package models

import (
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/cap_model"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
)

type DataFlange struct {
	Flange1    *flange_model.FlangeResult
	Flange2    *flange_model.FlangeResult
	Type1      flange_model.FlangeData_Type
	Type2      flange_model.FlangeData_Type
	Washer1    *flange_model.WasherResult
	Washer2    *flange_model.WasherResult
	Embed      *flange_model.EmbedResult
	Bolt       *flange_model.BoltResult
	Gasket     *flange_model.GasketResult
	TypeGasket string
	Dcp, B0    float64
}

type DataCap struct {
	Flange     *cap_model.FlangeResult
	Cap        *cap_model.CapResult
	FType      cap_model.FlangeData_Type
	CType      cap_model.CapData_Type
	Washer1    *cap_model.WasherResult
	Washer2    *cap_model.WasherResult
	Embed      *cap_model.EmbedResult
	Bolt       *cap_model.BoltResult
	Gasket     *cap_model.GasketResult
	TypeGasket string
	Dcp, B0    float64
}
