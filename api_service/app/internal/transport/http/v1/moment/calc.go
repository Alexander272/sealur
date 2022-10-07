package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model/cap_model"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model/flange_model"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model/float_model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initCalcRoutes(api *gin.RouterGroup) {
	calc := api.Group("/calc", h.middleware.UserIdentity)
	{
		calc.POST("/flange", h.calculateFlange)
		calc.POST("/cap", h.calculateCap)
		calc.POST("/float", h.calculateFloat)
	}
}

// @Summary Calculate Flange
// @Tags Sealur Moment -> calc-flange
// @Security ApiKeyAuth
// @Description расчет момента затяжки соединения фланец-фланец
// @ModuleID calculateFlange
// @Accept json
// @Produce json
// @Param data body flange_model.CalcFlange true "flange data"
// @Success 200 {object} models.DataResponse{data=calc_api.FlangeResponse}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/calc/flange/ [post]
func (h *Handler) calculateFlange(c *gin.Context) {
	var dto flange_model.CalcFlange
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	data, err := dto.NewFlange()
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	res, err := h.calcClient.CalculateFlange(c, data)
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: res})
}

// @Summary Calculate Cap
// @Tags Sealur Moment -> calc-cap
// @Security ApiKeyAuth
// @Description расчет момента затяжки соединения фланец-крышка
// @ModuleID calculateCap
// @Accept json
// @Produce json
// @Param data body cap_model.CalcCap true "cap data"
// @Success 200 {object} models.DataResponse{data=calc_api.CapResponse}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/calc/cap/ [post]
func (h *Handler) calculateCap(c *gin.Context) {
	var dto cap_model.CalcCap
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	data, err := dto.NewCap()
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	res, err := h.calcClient.CalculateCap(c, data)
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: res})
}

// @Summary Calculate Float
// @Tags Sealur Moment -> calc-float
// @Security ApiKeyAuth
// @Description расчет плавающей головки теплообменного аппарата
// @ModuleID calculateFloat
// @Accept json
// @Produce json
// @Param data body float_model.Calc true "float data"
// @Success 200 {object} models.DataResponse{data=calc_api.FloatResponse}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/calc/float/ [post]
func (h *Handler) calculateFloat(c *gin.Context) {
	var dto float_model.Calc
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	data, err := dto.New()
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	res, err := h.calcClient.CalculateFloat(c, data)
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: res})
}
