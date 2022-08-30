package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initCalcFlangeRoutes(api *gin.RouterGroup) {
	flange := api.Group("/calc/flange", h.middleware.UserIdentity)
	{
		flange.POST("/", h.calculate)
	}
}

// @Summary Calculate
// @Tags Sealur Moment -> calc-flange
// @Security ApiKeyAuth
// @Description расчет момента затяжки соединения фланец-фланец
// @ModuleID calculate
// @Accept json
// @Produce json
// @Param data body moment_model.CalcFlange true "size info"
// @Success 200 {object} models.DataResponse{data=moment_api.FlangeResponse}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/calc/flange/ [post]
func (h *Handler) calculate(c *gin.Context) {
	var dto moment_model.CalcFlange
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	// pressure, err := strconv.ParseFloat(dto.Pressure, 64)
	// if err != nil {
	// 	models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
	// 	return
	// }
	// axialForce, err := strconv.Atoi(dto.AxialForce)
	// if err != nil {
	// 	models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
	// 	return
	// }
	// bendingMoment, err := strconv.Atoi(dto.BendingMoment)
	// if err != nil {
	// 	models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
	// 	return
	// }
	// temp, err := strconv.ParseFloat(dto.Temp, 64)
	// if err != nil {
	// 	models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
	// 	return
	// }

	// flanges := moment_api.CalcFlangeRequest_Flanges_value[dto.Flanges]
	// typeB := moment_api.CalcFlangeRequest_Type_value[dto.TypeB]
	// condition := moment_api.CalcFlangeRequest_Condition_value[dto.Condition]
	// calculation := moment_api.CalcFlangeRequest_Calcutation_value[dto.Calculation]

	// flangesData := []*moment_api.FlangeData{}
	// flange1, err := h.newFlangeData(dto.FlangesData.First)
	// if err != nil {
	// 	models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
	// 	return
	// }
	// flangesData = append(flangesData, flange1)

	// if !dto.IsSameFlange {
	// 	flange2, err := h.newFlangeData(dto.FlangesData.Second)
	// 	if err != nil {
	// 		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
	// 		return
	// 	}
	// 	flangesData = append(flangesData, flange2)
	// }

	// bolts, err := dto.Bolts.NewBolts()
	// if err != nil {
	// 	models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
	// 	return
	// }

	// gasket, err := dto.Gasket.NewGasket()
	// if err != nil {
	// 	models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
	// 	return
	// }

	// var embed *moment_api.EmbedData
	// if dto.IsEmbedded {
	// 	embed, err = dto.Embed.NewEmbed()
	// 	if err != nil {
	// 		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
	// 		return
	// 	}
	// }

	// var washer []*moment_api.WasherData
	// if dto.IsUseWasher {
	// 	washer, err = dto.Washer.NewWasher()
	// 	if err != nil {
	// 		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
	// 		return
	// 	}
	// }

	// data := &moment_api.CalcFlangeRequest{
	// 	Pressure:       pressure,
	// 	AxialForce:     int32(axialForce),
	// 	BendingMoment:  int32(bendingMoment),
	// 	Temp:           temp,
	// 	IsWork:         dto.IsWork,
	// 	Flanges:        moment_api.CalcFlangeRequest_Flanges(flanges),
	// 	IsSameFlange:   dto.IsSameFlange,
	// 	IsEmbedded:     dto.IsEmbedded,
	// 	Type:           moment_api.CalcFlangeRequest_Type(typeB),
	// 	Condition:      moment_api.CalcFlangeRequest_Condition(condition),
	// 	Calculation:    moment_api.CalcFlangeRequest_Calcutation(calculation),
	// 	IsUseWasher:    dto.IsUseWasher,
	// 	IsNeedFormulas: dto.IsNeedFormulas,
	// 	FlangesData:    flangesData,
	// 	Bolts:          bolts,
	// 	Gasket:         gasket,
	// 	Washer:         washer,
	// 	Embed:          embed,
	// }

	data, err := dto.NewFlange()
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	res, err := h.calcFlangeClient.CalculateFlange(c, data)
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: res})
}

// func (h *Handler) newFlangeData(data moment_model.Flanges) (flange *moment_api.FlangeData, err error) {
// 	var size *moment_api.FlangeData_Size
// 	if data.StandartId == "another" {
// 		size, err = data.Size.NewSize()
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	var mat, rMat *moment_api.MaterialData
// 	if data.MarkId == "another" {
// 		mat, err = data.Material.NewMaterial()
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	if data.RingMarkId == "another" {
// 		rMat, err = data.RingMaterial.NewMaterial()
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	typeF := moment_api.FlangeData_Type_value[data.TypeF]
// 	corrosion, err := strconv.ParseFloat(data.Corrosion, 64)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var temp float64
// 	if data.Temp != "" {
// 		temp, err = strconv.ParseFloat(data.Temp, 64)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	flange = &moment_api.FlangeData{
// 		Type:         moment_api.FlangeData_Type(typeF),
// 		StandartId:   data.StandartId,
// 		MarkId:       data.MarkId,
// 		Material:     mat,
// 		Dy:           data.Dy,
// 		Py:           data.Py,
// 		Corrosion:    corrosion,
// 		Size:         size,
// 		Temp:         temp,
// 		RingMarkId:   data.RingMarkId,
// 		RingMaterial: rMat,
// 	}

// 	return flange, nil
// }
