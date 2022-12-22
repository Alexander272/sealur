package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment/material_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/material_model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initVoltageRoutes(api *gin.RouterGroup) {
	voltage := api.Group("/materials/voltage", h.middleware.UserIdentity, h.middleware.AccessForMomentAdmin)
	{
		voltage.POST("/", h.createVoltage)
		voltage.PUT("/:id", h.updateVoltage)
		voltage.DELETE("/:id", h.deleteVoltage)
	}
}

// @Summary Create Voltage
// @Tags Sealur Moment -> materials
// @Security ApiKeyAuth
// @Description создание данных для материала
// @ModuleID createVoltage
// @Accept json
// @Produce json
// @Param voltage body moment_model.CreateVoltageDTO true "voltage info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/voltage/ [post]
func (h *Handler) createVoltage(c *gin.Context) {
	var dto moment_model.CreateVoltageDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	voltage := make([]*material_model.Voltage, 0, len(dto.Voltage))
	for _, v := range dto.Voltage {
		voltage = append(voltage, &material_model.Voltage{
			Temperature: v.Temperature,
			Voltage:     v.Voltage,
		})
	}

	_, err := h.materialsClient.CreateVoltage(c, &material_api.CreateVoltageRequest{
		MarkId:  dto.MarkId,
		Voltage: voltage,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

// @Summary Update Voltage
// @Tags Sealur Moment -> materials
// @Security ApiKeyAuth
// @Description обновление данных для материала
// @ModuleID updateVoltage
// @Accept json
// @Produce json
// @Param id path string true "voltage id"
// @Param voltage body moment_model.UpdateVoltageDTO true "voltage info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/voltage/{id} [put]
func (h *Handler) updateVoltage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.UpdateVoltageDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.materialsClient.UpdateVoltage(c, &material_api.UpdateVoltageRequest{
		Id:          id,
		MarkId:      dto.MarkId,
		Temperature: dto.Temperature,
		Voltage:     dto.Voltage,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

// @Summary Delete Voltage
// @Tags Sealur Moment -> materials
// @Security ApiKeyAuth
// @Description Удаление данных для материала
// @ModuleID deleteVoltage
// @Accept json
// @Produce json
// @Param id path string true "voltage id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/voltage/{id} [delete]
func (h *Handler) deleteVoltage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.materialsClient.DeleteVoltage(c, &material_api.DeleteVoltageRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
