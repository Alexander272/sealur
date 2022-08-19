package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initGasketDataRoutes(api *gin.RouterGroup) {
	gasket := api.Group("/gasket-data", h.middleware.UserIdentity, h.middleware.AccessForMomentAdmin)
	{
		gasket.POST("/", h.createGasketData)
		gasket.POST("/many", h.createManyGasketData)
		gasket.PUT("/:id", h.updateGasketData)
		gasket.PUT("/typeId/:gasketId", h.updateGasketTypeId)
		gasket.DELETE("/:id", h.deleteGasketData)
	}
}

// @Summary Create Many Gasket Data
// @Tags Sealur Moment -> gasket-data
// @Security ApiKeyAuth
// @Description создание нескольких данных для прокладки
// @ModuleID createManyGasketData
// @Accept json
// @Produce json
// @Param gasketData body moment_model.ManyGasketDataDTO true "gasket data info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/gasket-data/many [post]
func (h *Handler) createManyGasketData(c *gin.Context) {
	var dto moment_model.ManyGasketDataDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	data := []*moment_api.CreateManyGasketDataRequest_Data{}
	for _, mgdd := range dto.Data {
		data = append(data, &moment_api.CreateManyGasketDataRequest_Data{
			PermissiblePres: mgdd.PermissiblePres,
			Compression:     mgdd.Compression,
			Epsilon:         mgdd.Epsilon,
			Thickness:       mgdd.Thickness,
		})
	}

	_, err := h.gasketClient.CreateManyGasketData(c, &moment_api.CreateManyGasketDataRequest{
		GasketId: dto.GasketId,
		TypeId:   dto.TypeId,
		Data:     data,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

// @Summary Create Gasket Data
// @Tags Sealur Moment -> gasket-data
// @Security ApiKeyAuth
// @Description создание данных для прокладки
// @ModuleID createGasketData
// @Accept json
// @Produce json
// @Param gasketData body moment_model.GasketDataDTO true "gasket data info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/gasket-data/ [post]
func (h *Handler) createGasketData(c *gin.Context) {
	var dto moment_model.GasketDataDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.gasketClient.CreateGasketData(c, &moment_api.CreateGasketDataRequest{
		GasketId:        dto.GasketId,
		PermissiblePres: dto.PermissiblePres,
		Compression:     dto.Compression,
		Epsilon:         dto.Epsilon,
		Thickness:       dto.Thickness,
		TypeId:          dto.TypeId,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

// @Summary Update Gasket Data
// @Tags Sealur Moment -> gasket-data
// @Security ApiKeyAuth
// @Description обновление данных для прокладки
// @ModuleID updateGasketData
// @Accept json
// @Produce json
// @Param id path string true "gasket data id"
// @Param gasketData body moment_model.GasketDataDTO true "gasket data info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/gasket-data/{id} [put]
func (h *Handler) updateGasketData(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.GasketDataDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.gasketClient.UpdateGasketData(c, &moment_api.UpdateGasketDataRequest{
		Id:              id,
		GasketId:        dto.GasketId,
		PermissiblePres: dto.PermissiblePres,
		Compression:     dto.Compression,
		Epsilon:         dto.Epsilon,
		Thickness:       dto.Thickness,
		TypeId:          dto.TypeId,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

// @Summary Update Gasket Data
// @Tags Sealur Moment -> gasket-data
// @Security ApiKeyAuth
// @Description обновление данных для прокладки
// @ModuleID updateGasketData
// @Accept json
// @Produce json
// @Param gasketId path string true "gasket data id"
// @Param gasketData body moment_model.GasketDataDTO true "gasket data info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/gasket-data/typeId/{gasketId} [put]
func (h *Handler) updateGasketTypeId(c *gin.Context) {
	gasketId := c.Param("gasketId")
	if gasketId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.GasketDataDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.gasketClient.UpdateGasketTypeId(c, &moment_api.UpdateGasketTypeIdRequest{
		GasketId: gasketId,
		TypeId:   dto.TypeId,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated"})
}

// @Summary Delete Gasket Data
// @Tags Sealur Moment -> gasket-data
// @Security ApiKeyAuth
// @Description Удаление данных для прокладки
// @ModuleID deleteGasketData
// @Accept json
// @Produce json
// @Param id path string true "gasket data id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/gasket-data/{id} [delete]
func (h *Handler) deleteGasketData(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.gasketClient.DeleteGasketData(c, &moment_api.DeleteGasketDataRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
