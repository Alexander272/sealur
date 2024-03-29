package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initFlangeRoutes(api *gin.RouterGroup) {
	flange := api.Group("/flange-sizes", h.middleware.UserIdentity, h.middleware.AccessForMomentAdmin)
	{
		flange.GET("/", h.getFlangeSize)
		flange.POST("/", h.createFlangeSize)
		flange.POST("/several", h.createFlangeSizes)
		flange.PATCH("/:id", h.updateFlangeSize)
		flange.DELETE("/:id", h.deleteFlangeSize)
	}
}

// @Summary Get Flange Size
// @Tags Sealur Moment -> flange-sizes
// @Security ApiKeyAuth
// @Description создание размеров
// @ModuleID getFlangeSize
// @Accept json
// @Produce json
// @Param standartId query string true "standart Id"
// @Success 200 {object} models.DataResponse{data=flange_api.FullFlangeSizeResponse}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/flange-sizes/ [get]
func (h *Handler) getFlangeSize(c *gin.Context) {
	id := c.Query("standartId")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	sizes, err := h.flangeClient.GetFlangeSize(c, &flange_api.GetFullFlangeSizeRequest{
		StandId: id,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: sizes})
}

// @Summary Create Flange Size
// @Tags Sealur Moment -> flange-sizes
// @Security ApiKeyAuth
// @Description создание размеров
// @ModuleID createFlangeSize
// @Accept json
// @Produce json
// @Param size body moment_model.FlangeSizeDTO true "size info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/flange-sizes/ [post]
func (h *Handler) createFlangeSize(c *gin.Context) {
	var dto moment_model.FlangeSizeDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.flangeClient.CreateFlangeSize(c, &flange_api.CreateFlangeSizeRequest{
		StandId: dto.StandId,
		Pn:      dto.Pn,
		Dn:      dto.Dn,
		Dmm:     dto.Dmm,
		D:       dto.D,
		D6:      dto.D6,
		DOut:    dto.DOut,
		X:       dto.X,
		A:       dto.A,
		H:       dto.H,
		S0:      dto.S0,
		S1:      dto.S1,
		Length:  dto.Length,
		Count:   dto.Count,
		BoltId:  dto.BoltId,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

// @Summary Create Flange Sizes
// @Tags Sealur Moment -> flange-sizes
// @Security ApiKeyAuth
// @Description создание размеров
// @ModuleID createFlangeSizes
// @Accept json
// @Produce json
// @Param size body []moment_model.FlangeSizeDTO true "size info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/flange-sizes/several [post]
func (h *Handler) createFlangeSizes(c *gin.Context) {
	var dto []moment_model.FlangeSizeDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	sizes := []*flange_api.CreateFlangeSizeRequest{}
	for _, sd := range dto {
		sizes = append(sizes, &flange_api.CreateFlangeSizeRequest{
			StandId: sd.StandId,
			Pn:      sd.Pn,
			Dn:      sd.Dn,
			Dmm:     sd.Dmm,
			D:       sd.D,
			D6:      sd.D6,
			DOut:    sd.DOut,
			X:       sd.X,
			A:       sd.A,
			H:       sd.H,
			S0:      sd.S0,
			S1:      sd.S1,
			Length:  sd.Length,
			Count:   sd.Count,
			BoltId:  sd.BoltId,
			Row:     sd.Row,
		})
	}

	_, err := h.flangeClient.CreateFlangeSizes(c, &flange_api.CreateFlangeSizesRequest{Sizes: sizes})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

// @Summary Update Flange Size
// @Tags Sealur Moment -> flange-sizes
// @Security ApiKeyAuth
// @Description обновление размеров
// @ModuleID updateFlangeSize
// @Accept json
// @Produce json
// @Param id path string true "size id"
// @Param size body moment_model.FlangeSizeDTO true "size info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/flange-sizes/{id} [put]
func (h *Handler) updateFlangeSize(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.FlangeSizeDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.flangeClient.UpdateFlangeSize(c, &flange_api.UpdateFlangeSizeRequest{
		Id:      id,
		StandId: dto.StandId,
		Pn:      dto.Pn,
		Dn:      dto.Dn,
		Dmm:     dto.Dmm,
		D:       dto.D,
		D6:      dto.D6,
		DOut:    dto.DOut,
		X:       dto.X,
		A:       dto.A,
		H:       dto.H,
		S0:      dto.S0,
		S1:      dto.S1,
		Length:  dto.Length,
		Count:   dto.Count,
		BoltId:  dto.BoltId,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

// @Summary Delete Flange Size
// @Tags Sealur Moment -> flange-sizes
// @Security ApiKeyAuth
// @Description Удаление болта
// @ModuleID deleteFlangeSize
// @Accept json
// @Produce json
// @Param id path string true "size id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/flange-sizes/{id} [delete]
func (h *Handler) deleteFlangeSize(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.flangeClient.DeleteFlangeSize(c, &flange_api.DeleteFlangeSizeRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
