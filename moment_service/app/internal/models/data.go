package models

import "github.com/Alexander272/sealur_proto/api/moment_api"

type DataFlange struct {
	Flange1    *moment_api.FlangeResult
	Flange2    *moment_api.FlangeResult
	Type1      moment_api.FlangeData_Type
	Type2      moment_api.FlangeData_Type
	Washer1    *moment_api.WasherResult
	Washer2    *moment_api.WasherResult
	Embed      *moment_api.EmbedResult
	Bolt       *moment_api.BoltResult
	Gasket     *moment_api.GasketResult
	TypeGasket string
	Dcp, B0    float64
}

type DataCap struct {
	Flange     *moment_api.FlangeResult
	Cap        *moment_api.CapResult
	FType      moment_api.FlangeData_Type
	CType      moment_api.CapData_Type
	Washer1    *moment_api.WasherResult
	Washer2    *moment_api.WasherResult
	Embed      *moment_api.EmbedResult
	Bolt       *moment_api.BoltResult
	Gasket     *moment_api.GasketResult
	TypeGasket string
	Dcp, B0    float64
}
