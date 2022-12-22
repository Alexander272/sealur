package models

import (
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/cap_model"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/dev_cooling_model"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/express_circle_model"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/express_rectangle_model"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/float_model"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/gas_cooling_model"
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
	TypeGasket flange_model.GasketData_Type
}

type DataFlangeOld struct {
	Flange1    *flange_model.OldFlangeResult
	Flange2    *flange_model.OldFlangeResult
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

type DataFloat struct {
	Flange     *float_model.FlangeResult
	Cap        *float_model.CapResult
	Bolt       *float_model.BoltResult
	Gasket     *float_model.GasketResult
	TypeGasket string
	Dcp, B0    float64
}

type DataDevCooling struct {
	Cap        *dev_cooling_model.CapResult
	TubeSheet  *dev_cooling_model.TubeSheetResult
	Tube       *dev_cooling_model.TubeResult
	Bolt       *dev_cooling_model.BoltResult
	Gasket     *dev_cooling_model.GasketResult
	TypeGasket dev_cooling_model.GasketData_Type
}

type DataGasCooling struct {
	Data       *gas_cooling_model.DataResult
	Bolt       *gas_cooling_model.BoltResult
	Gasket     *gas_cooling_model.GasketResult
	TypeGasket gas_cooling_model.GasketData_Type
}

type DataExCircle struct {
	Bolt       *express_circle_model.BoltResult
	Gasket     *express_circle_model.GasketResult
	TypeGasket express_circle_model.GasketData_Type
}

type DataExRect struct {
	Bolt       *express_rectangle_model.BoltResult
	Gasket     *express_rectangle_model.GasketResult
	TypeGasket express_rectangle_model.GasketData_Type
}
