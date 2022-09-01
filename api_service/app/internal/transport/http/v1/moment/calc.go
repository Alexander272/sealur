package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initCalcRoutes(api *gin.RouterGroup) {
	calc := api.Group("/calc", h.middleware.UserIdentity)
	{
		calc.POST("/flange", h.calculateFlange)
		calc.POST("/cap", h.calculateCap)
	}
}

// @Summary Calculate Flange
// @Tags Sealur Moment -> calc-flange
// @Security ApiKeyAuth
// @Description расчет момента затяжки соединения фланец-фланец
// @ModuleID calculateFlange
// @Accept json
// @Produce json
// @Param data body moment_model.CalcFlange true "flange data"
// @Success 200 {object} models.DataResponse{data=moment_api.FlangeResponse}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/calc/flange/ [post]
func (h *Handler) calculateFlange(c *gin.Context) {
	var dto moment_model.CalcFlange
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

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

// @Summary Calculate Cap
// @Tags Sealur Moment -> calc-cap
// @Security ApiKeyAuth
// @Description расчет момента затяжки соединения фланец-крышка
// @ModuleID calculateCap
// @Accept json
// @Produce json
// @Param data body moment_model.CalcCap true "cap data"
// @Success 200 {object} models.DataResponse{data=moment_api.CapData}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/calc/cap/ [post]
func (h *Handler) calculateCap(c *gin.Context) {
	var dto moment_model.CalcCap
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	data, err := dto.NewCap()
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	res, err := h.calcCapClient.CalculateCap(c, data)
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: res})
}
