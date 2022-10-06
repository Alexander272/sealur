package moment

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initTypeGasketRoutes(api *gin.RouterGroup) {
	gasket := api.Group("/type-gasket", h.middleware.UserIdentity)
	{
		gasket.GET("/", h.getTypeGasket)
		gasket = gasket.Group("/", h.middleware.AccessForMomentAdmin)
		{
			gasket.POST("/", h.createTypeGasket)
			gasket.PATCH("/:id", h.updateTypeGasket)
			gasket.DELETE("/:id", h.deleteTypeGasket)
		}
	}
}

// @Summary Get Type Gasket
// @Tags Sealur Moment -> type-gasket
// @Security ApiKeyAuth
// @Description получение типов прокладок
// @ModuleID getTypeGasket
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{Data=[]gasket_api.GasketType}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/type-gasket/ [get]
func (h *Handler) getTypeGasket(c *gin.Context) {
	gasket, err := h.gasketClient.GetGasketType(c, &gasket_api.GetGasketTypeRequest{})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: gasket.GasketType, Count: len(gasket.GasketType)})
}

// @Summary Create Type Gasket
// @Tags Sealur Moment -> type-gasket
// @Security ApiKeyAuth
// @Description создание типа прокладки
// @ModuleID createTypeGasket
// @Accept json
// @Produce json
// @Param typeGasket body moment_model.TypeGasketDTO true "type-gasket info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/type-gasket/ [post]
func (h *Handler) createTypeGasket(c *gin.Context) {
	var dto moment_model.TypeGasketDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	gasket, err := h.gasketClient.CreateGasketType(c, &gasket_api.CreateGasketTypeRequest{Title: dto.Title})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-moment/type-gasket/%s", gasket.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: gasket.Id, Message: "Created"})
}

// @Summary Update Type Gasket
// @Tags Sealur Moment -> type-gasket
// @Security ApiKeyAuth
// @Description обновление типа прокладки
// @ModuleID updateTypeGasket
// @Accept json
// @Produce json
// @Param id path string true "type gasket id"
// @Param typeGasket body moment_model.TypeGasketDTO true "type gasket info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/type-gasket/{id} [put]
func (h *Handler) updateTypeGasket(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.GasketDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.gasketClient.UpdateGasketType(c, &gasket_api.UpdateGasketTypeRequest{Id: id, Title: dto.Title})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

// @Summary Delete Type Gasket
// @Tags Sealur Moment -> type-gasket
// @Security ApiKeyAuth
// @Description Удаление типа прокладки
// @ModuleID deleteTypeGasket
// @Accept json
// @Produce json
// @Param id path string true "type gasket id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/type-gasket/{id} [delete]
func (h *Handler) deleteTypeGasket(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.gasketClient.DeleteGasketType(c, &gasket_api.DeleteGasketTypeRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
