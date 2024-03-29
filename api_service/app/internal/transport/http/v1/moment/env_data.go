package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initEnvDataRoutes(api *gin.RouterGroup) {
	env := api.Group("/env-data", h.middleware.UserIdentity, h.middleware.AccessForMomentAdmin)
	{
		env.POST("/", h.createEnvData)
		env.POST("/many", h.createManyEnvData)
		env.PUT("/:id", h.updateEnvData)
		env.DELETE("/:id", h.deleteEnvData)
	}
}

// @Summary Create Env Data
// @Tags Sealur Moment -> env-data
// @Security ApiKeyAuth
// @Description создание данных среды
// @ModuleID createEnvData
// @Accept json
// @Produce json
// @Param env body moment_model.EnvDataDTO true "env data info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/env-data/ [post]
func (h *Handler) createEnvData(c *gin.Context) {
	var dto moment_model.EnvDataDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.gasketClient.CreateEnvData(c, &gasket_api.CreateEnvDataRequest{
		EnvId:        dto.EnvId,
		GasketId:     dto.GasketId,
		M:            dto.M,
		SpecificPres: dto.SpecificPres,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

// @Summary Create Env Data
// @Tags Sealur Moment -> env-data
// @Security ApiKeyAuth
// @Description создание данных среды
// @ModuleID createEnvData
// @Accept json
// @Produce json
// @Param env body moment_model.EnvDataDTO true "env data info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/env-data/many [post]
func (h *Handler) createManyEnvData(c *gin.Context) {
	var dto moment_model.ManyEnvDataDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	data := []*gasket_api.CreateManyEnvDataRequest_Data{}
	for _, mgdd := range dto.Data {
		data = append(data, &gasket_api.CreateManyEnvDataRequest_Data{
			SpecificPres: mgdd.SpecificPres,
			M:            mgdd.M,
			EnvId:        mgdd.EnvId,
		})
	}

	_, err := h.gasketClient.CreateManyEnvData(c, &gasket_api.CreateManyEnvDataRequest{
		GasketId: dto.GasketId,
		Data:     data,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

// @Summary Update Env Data
// @Tags Sealur Moment -> env-data
// @Security ApiKeyAuth
// @Description обновление данных среды
// @ModuleID updateEnv
// @Accept json
// @Produce json
// @Param id path string true "env data id"
// @Param env body moment_model.EnvDataDTO true "env data info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/env-data/{id} [put]
func (h *Handler) updateEnvData(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.EnvDataDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.gasketClient.UpdateEnvData(c, &gasket_api.UpdateEnvDataRequest{
		Id:           id,
		EnvId:        dto.EnvId,
		GasketId:     dto.GasketId,
		M:            dto.M,
		SpecificPres: dto.SpecificPres,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

// @Summary Delete Env Data
// @Tags Sealur Moment -> env-data
// @Security ApiKeyAuth
// @Description Удаление данных среды
// @ModuleID deleteEnv
// @Accept json
// @Produce json
// @Param id path string true "env data id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/env-data/{id} [delete]
func (h *Handler) deleteEnvData(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.gasketClient.DeleteEnvData(c, &gasket_api.DeleteEnvDataRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
