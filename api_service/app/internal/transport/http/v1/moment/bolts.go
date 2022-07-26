package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initBoltsRoutes(api *gin.RouterGroup) {
	bolts := api.Group("/bolts", h.middleware.UserIdentity)
	{
		bolts.GET("/", h.getBolts)
		bolts = bolts.Group("/", h.middleware.AccessForMomentAdmin)
		{
			bolts.POST("/", h.createBolt)
			bolts.PATCH("/:id", h.updateBolt)
			bolts.DELETE("/:id", h.deleteBolt)
		}
	}
}

// @Summary Get Bolts
// @Tags Sealur Moment -> bolts
// @Security ApiKeyAuth
// @Description получение болтов
// @ModuleID getBolts
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{Data=[]moment_api.Bolt}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/bolts/ [get]
func (h *Handler) getBolts(c *gin.Context) {
	bolts, err := h.flangeClient.GetBolts(c, &moment_api.GetBoltsRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: bolts.Bolts, Count: len(bolts.Bolts)})
}

// @Summary Create Bolt
// @Tags Sealur Moment -> bolts
// @Security ApiKeyAuth
// @Description создание болта
// @ModuleID createBolt
// @Accept json
// @Produce json
// @Param bolt body moment_model.BoltDTO true "bolt info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/bolts/ [post]
func (h *Handler) createBolt(c *gin.Context) {
	var dto moment_model.BoltDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.flangeClient.CreateBolt(c, &moment_api.CreateBoltRequest{
		Title:    dto.Title,
		Diameter: dto.Diameter,
		Area:     dto.Area,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

// @Summary Update Bolt
// @Tags Sealur Moment -> bolts
// @Security ApiKeyAuth
// @Description обновление среды
// @ModuleID updateBolt
// @Accept json
// @Produce json
// @Param id path string true "bolt id"
// @Param bolt body moment_model.BoltDTO true "bolt info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/bolts/{id} [put]
func (h *Handler) updateBolt(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.BoltDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.flangeClient.UpdateBolt(c, &moment_api.UpdateBoltRequest{
		Id:       id,
		Title:    dto.Title,
		Diameter: dto.Diameter,
		Area:     dto.Area,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

// @Summary Delete Bolt
// @Tags Sealur Moment -> bolts
// @Security ApiKeyAuth
// @Description Удаление болта
// @ModuleID deleteBolt
// @Accept json
// @Produce json
// @Param id path string true "bolt id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/bolts/{id} [delete]
func (h *Handler) deleteBolt(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.flangeClient.DeleteBolt(c, &moment_api.DeleteBoltRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
