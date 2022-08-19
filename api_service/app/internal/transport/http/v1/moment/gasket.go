package moment

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initGasketRoutes(api *gin.RouterGroup) {
	gasket := api.Group("/gasket", h.middleware.UserIdentity)
	{
		gasket.GET("/", h.getGasket)
		gasket.GET("/full-data", h.getFullData)
		gasket = gasket.Group("/", h.middleware.AccessForMomentAdmin)
		{
			gasket.POST("/", h.createGasket)
			gasket.PUT("/:id", h.updateGasket)
			gasket.DELETE("/:id", h.deleteGasket)
		}
	}
}

// @Summary Get Gasket
// @Tags Sealur Moment -> gasket
// @Security ApiKeyAuth
// @Description получение всех типов прокладок
// @ModuleID getGasket
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{Data=[]moment_api.Gasket}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/gasket/ [get]
func (h *Handler) getGasket(c *gin.Context) {
	gasket, err := h.gasketClient.GetGasket(c, &moment_api.GetGasketRequest{})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: gasket.Gasket, Count: len(gasket.Gasket)})
}

// @Summary Get Full Data
// @Tags Sealur Moment -> gasket
// @Security ApiKeyAuth
// @Description получение данных для прокладки
// @ModuleID getFullData
// @Accept json
// @Produce json
// @Param gasketId query string true "gasket id"
// @Success 200 {object} models.DataResponse{Data=[]moment_api.Gasket}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/gasket/full-data [get]
func (h *Handler) getFullData(c *gin.Context) {
	id := c.Query("gasketId")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	data, err := h.gasketClient.GetFullData(c, &moment_api.GetFullDataRequest{GasketId: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: data})
}

// @Summary Create Gasket
// @Tags Sealur Moment -> gasket
// @Security ApiKeyAuth
// @Description создание прокладки
// @ModuleID createGasket
// @Accept json
// @Produce json
// @Param gasket body moment_model.GasketDTO true "gasket info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/gasket/ [post]
func (h *Handler) createGasket(c *gin.Context) {
	var dto moment_model.GasketDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	gasket, err := h.gasketClient.CreateGasket(c, &moment_api.CreateGasketRequest{Title: dto.Title})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-moment/gasket/%s", gasket.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: gasket.Id, Message: "Created"})
}

// @Summary Update Gasket
// @Tags Sealur Moment -> gasket
// @Security ApiKeyAuth
// @Description обновление прокладки
// @ModuleID updateGasket
// @Accept json
// @Produce json
// @Param id path string true "gasket id"
// @Param gasket body moment_model.GasketDTO true "gasket info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/gasket/{id} [put]
func (h *Handler) updateGasket(c *gin.Context) {
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

	_, err := h.gasketClient.UpdateGasket(c, &moment_api.UpdateGasketRequest{Id: id, Title: dto.Title})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

// @Summary Delete Gasket
// @Tags Sealur Moment -> gasket
// @Security ApiKeyAuth
// @Description Удаление прокладки
// @ModuleID deleteGasket
// @Accept json
// @Produce json
// @Param id path string true "gasket id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/gasket/{id} [delete]
func (h *Handler) deleteGasket(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.gasketClient.DeleteGasket(c, &moment_api.DeleteGasketRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
