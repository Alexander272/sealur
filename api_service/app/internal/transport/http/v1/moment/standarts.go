package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto/moment_proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initStandartsRoutes(api *gin.RouterGroup) {
	standarts := api.Group("/standarts", h.middleware.UserIdentity)
	{
		standarts.GET("/", h.getStandarts)
		standarts = standarts.Group("/", h.middleware.AccessForMomentAdmin)
		{
			standarts.POST("/", h.createStandart)
			standarts.PATCH("/:id", h.updateStandart)
			standarts.DELETE("/:id", h.deleteStandart)
		}
	}
}

// @Summary Get Standarts
// @Tags Sealur Moment -> standarts
// @Security ApiKeyAuth
// @Description получение стандартов
// @ModuleID getStandarts
// @Accept json
// @Produce json
// @Param typeId query string true "type id"
// @Success 200 {object} models.DataResponse{Data=[]moment_proto.Standart}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/standarts/ [get]
func (h *Handler) getStandarts(c *gin.Context) {
	typeId := c.Query("typeId")
	if typeId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty typeId param")
		return
	}

	standarts, err := h.flangeClient.GetStandarts(c, &moment_proto.GetStandartsRequest{TypeId: typeId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: standarts.Standarts, Count: len(standarts.Standarts)})
}

// @Summary Create Standart
// @Tags Sealur Moment -> standarts
// @Security ApiKeyAuth
// @Description создание болта
// @ModuleID createStandart
// @Accept json
// @Produce json
// @Param standart body models.MomentStandartDTO true "standart info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/standarts/ [post]
func (h *Handler) createStandart(c *gin.Context) {
	var dto models.MomentStandartDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	stand, err := h.flangeClient.CreateStandart(c, &moment_proto.CreateStandartRequest{
		Title:  dto.Title,
		TypeId: dto.TypeId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Id: stand.Id, Message: "Created"})
}

// @Summary Update Standart
// @Tags Sealur Moment -> standarts
// @Security ApiKeyAuth
// @Description обновление среды
// @ModuleID updateStandart
// @Accept json
// @Produce json
// @Param id path string true "standart id"
// @Param standart body models.MomentStandartDTO true "standart info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/standarts/{id} [put]
func (h *Handler) updateStandart(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto models.MomentStandartDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.flangeClient.UpdateStandart(c, &moment_proto.UpdateStandartRequest{
		Id:     id,
		Title:  dto.Title,
		TypeId: dto.TypeId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

// @Summary Delete Standart
// @Tags Sealur Moment -> standarts
// @Security ApiKeyAuth
// @Description Удаление болта
// @ModuleID deleteStandart
// @Accept json
// @Produce json
// @Param id path string true "standart id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/standarts/{id} [delete]
func (h *Handler) deleteStandart(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.flangeClient.DeleteStandart(c, &moment_proto.DeleteStandartRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
