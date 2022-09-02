package formulas

import (
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type Flange interface {
	GetFormulasForFlange(
		TypeGasket, Condition string, IsWork, IsUseWasher, IsEmbedded bool, data models.DataFlange, result moment_api.FlangeResponse,
		calculation moment_api.CalcFlangeRequest_Calcutation, gamma_, yb_, yp_ float64,
	) *moment_api.CalcFlangeFormulas
}

type Cap interface {
	GetFormulasForCap(
		TypeGasket, Condition string, IsWork, IsUseWasher, IsEmbedded bool, data models.DataCap, result moment_api.CapResponse,
		calculation moment_api.CalcCapRequest_Calcutation, gamma_, yb_, yp_ float64,
	) *moment_api.CalcCapFormulas
}

type FormulasService struct {
	Flange
	Cap
}

func NewFormulasService() *FormulasService {
	return &FormulasService{
		Flange: NewFlangeService(),
		Cap:    NewCapService(),
	}
}
