package pro

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initStFlRoutes(api *gin.RouterGroup) {
	st := api.Group("st-fl")
	{
		st.GET("/", h.getStFl)
		st.POST("/", h.createStFl)
		st.PUT("/:id", h.updateStFl)
		st.DELETE("/:id", h.deleteStFl)
	}
}

// @Summary Get Stand/Flange
// @Tags Sealur Pro -> st-fl
// @Security ApiKeyAuth
// @Description получение списка стандартов (только для снп)
// @ModuleID getStFl
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{data=[]proto.StFl}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/st-fl [get]
func (h *Handler) getStFl(c *gin.Context) {
	st, err := h.proClient.GetStFl(c, &proto.GetStFlRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: st.Stfl})
}

// @Summary Create Stand/Flange
// @Tags Sealur Pro -> st-fl
// @Security ApiKeyAuth
// @Description создание элемента списка стандартов (только для снп)
// @ModuleID createStFl
// @Accept json
// @Produce json
// @Param data body models.StFlDTO true "st/fl info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/st-fl [post]
func (h *Handler) createStFl(c *gin.Context) {
	var dto models.StFlDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	request := &proto.CreateStFlRequest{
		StandId:  dto.StandId,
		FlangeId: dto.FlangeId,
	}

	st, err := h.proClient.CreateStFl(c, request)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/st-fl/%s", st.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: st.Id, Message: "Created"})
}

// @Summary Update Stand/Flange
// @Tags Sealur Pro -> st-fl
// @Security ApiKeyAuth
// @Description обновление элемента списка стандартов (только для снп)
// @ModuleID updateStFl
// @Accept json
// @Produce json
// @Param id path string true "st/fl id"
// @Param data body models.StFlDTO true "st/fl info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/st-fl/{id} [put]
func (h *Handler) updateStFl(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto models.StFlDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	request := &proto.UpdateStFlRequest{
		StandId:  dto.StandId,
		FlangeId: dto.FlangeId,
	}

	st, err := h.proClient.UpdateStFl(c, request)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: st.Id, Message: "Updated"})
}

// @Summary Delete Stand/Flange
// @Tags Sealur Pro -> st-fl
// @Security ApiKeyAuth
// @Description удаление элемента списка стандартов (только для снп)
// @ModuleID deleteStFl
// @Accept json
// @Produce json
// @Param id path string true "st/fl id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/st-fl/{id} [delete]
func (h *Handler) deleteStFl(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	st, err := h.proClient.DeleteStFl(c, &proto.DeleteStFlRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: st.Id, Message: "Deleted"})
}
