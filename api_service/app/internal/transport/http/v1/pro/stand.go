package pro

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initStandRoutes(api *gin.RouterGroup) {
	stand := api.Group("/standards")
	{
		stand.GET("/", h.getStands)
		stand.POST("/", h.createStand)
		stand.PUT("/:id", h.updateStand)
		stand.DELETE("/:id", h.deleteStand)
	}
}

// @Summary Get Stand
// @Tags Sealur Pro -> standards
// @Security ApiKeyAuth
// @Description получение всех стандартов на прокладки
// @ModuleID getStands
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{data=[]proto.Stand}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/standards [get]
func (h *Handler) getStands(c *gin.Context) {
	st, err := h.proClient.GetAllStands(c, &proto.GetStandsRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: st.Stands})
}

// @Summary Create Stand
// @Tags Sealur Pro -> standards
// @Security ApiKeyAuth
// @Description создание стандарта на прокладки
// @ModuleID createStand
// @Accept json
// @Produce json
// @Param data body models.StandDTO true "standard info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/standards [post]
func (h *Handler) createStand(c *gin.Context) {
	var dto models.StandDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data sent")
		return
	}

	id, err := h.proClient.CreateStand(c, &proto.CreateStandRequest{Title: dto.Title})
	if err != nil {
		if errors.Is(err, models.ErrStandAlreadyExists) {
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
			return
		}
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/standards/%s", id.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: id.Id, Message: "Created"})
}

// @Summary Update Stand
// @Tags Sealur Pro -> standards
// @Security ApiKeyAuth
// @Description обновление стандарта на прокладки
// @ModuleID updateStand
// @Accept json
// @Produce json
// @Param data body models.StandDTO true "standard info"
// @Param id path string true "standard id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/standards/{id} [put]
func (h *Handler) updateStand(c *gin.Context) {
	var dto models.StandDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}
	stId := c.Param("id")
	if stId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	st, err := h.proClient.UpdateStand(c, &proto.UpdateStandRequest{Id: stId, Title: dto.Title})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: st.Id, Message: "Updated"})
}

// @Summary Delete Stand
// @Tags Sealur Pro -> standards
// @Security ApiKeyAuth
// @Description удаление стандарта на прокладки
// @ModuleID deleteStand
// @Accept json
// @Produce json
// @Param id path string true "standard id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/standards/{id} [delete]
func (h *Handler) deleteStand(c *gin.Context) {
	stId := c.Param("id")
	if stId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	st, err := h.proClient.DeleteStand(c, &proto.DeleteStandRequest{Id: stId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: st.Id, Message: "Deleted"})
}
