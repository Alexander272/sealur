package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto/moment_proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initAlphaRoutes(api *gin.RouterGroup) {
	alpha := api.Group("/materials/alpha", h.middleware.UserIdentity, h.middleware.AccessForMomentAdmin)
	{
		alpha.POST("/", h.createAlpha)
		alpha.PUT("/:id", h.updateAlpha)
		alpha.DELETE("/:id", h.deleteAlpha)
	}
}

// @Summary Create Alpha
// @Tags Sealur Moment -> materials
// @Security ApiKeyAuth
// @Description создание данных для материала
// @ModuleID createAlpha
// @Accept json
// @Produce json
// @Param alpha body models.CreateAlphaDTO true "alpha info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/alpha/ [post]
func (h *Handler) createAlpha(c *gin.Context) {
	var dto models.CreateAlphaDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	alpha := make([]*moment_proto.Alpha, 0, len(dto.Alpha))
	for _, v := range dto.Alpha {
		item := moment_proto.Alpha(v)
		alpha = append(alpha, &item)
	}

	_, err := h.materialsClient.CreateAlpha(c, &moment_proto.CreateAlphaRequest{
		MarkId: dto.MarkId,
		Alpha:  alpha,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

// @Summary Update Alpha
// @Tags Sealur Moment -> materials
// @Security ApiKeyAuth
// @Description обновление данных для материала
// @ModuleID updateAlpha
// @Accept json
// @Produce json
// @Param id path string true "alpha id"
// @Param alpha body models.UpdateAlphaDTO true "alpha info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/alpha/{id} [put]
func (h *Handler) updateAlpha(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto models.UpdateAlphaDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.materialsClient.UpdateAlpha(c, &moment_proto.UpdateAlphaRequest{
		Id:          id,
		MarkId:      dto.MarkId,
		Temperature: dto.Temperature,
		Alpha:       dto.Alpha,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

// @Summary Delete Alpha
// @Tags Sealur Moment -> materials
// @Security ApiKeyAuth
// @Description Удаление данных для материала
// @ModuleID deleteAlpha
// @Accept json
// @Produce json
// @Param id path string true "alpha id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/alpha/{id} [delete]
func (h *Handler) deleteAlpha(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.materialsClient.DeleteAlpha(c, &moment_proto.DeleteAlphaRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
