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
// @Param size body moment_model.FlangeSizeDTO true "size info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/flange-sizes/ [post]
func (h *Handler) calculate(c *gin.Context) {
	var dto moment_model.FlangeSizeDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	// _, err := h.calcFlangeClient.CalculateFlange(c, &moment_api.{
	// 	StandId: dto.StandId,
	// 	Pn:      dto.Pn,
	// 	D:       dto.D,
	// 	D6:      dto.D6,
	// 	DOut:    dto.DOut,
	// 	H:       dto.H,
	// 	S0:      dto.S0,
	// 	S1:      dto.S1,
	// 	Length:  dto.Length,
	// 	Count:   dto.Count,
	// 	BoltId:  dto.BoltId,
	// })
	// if err != nil {
	// 	models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
	// 	return
	// }

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}
