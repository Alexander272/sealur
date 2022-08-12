package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initElasticityRoutes(api *gin.RouterGroup) {
	elasticity := api.Group("/materials/elasticity", h.middleware.UserIdentity, h.middleware.AccessForMomentAdmin)
	{
		elasticity.POST("/", h.createElasticity)
		elasticity.PUT("/:id", h.updateElasticity)
		elasticity.DELETE("/:id", h.deleteElasticity)
	}
}

// @Summary Create Elasticity
// @Tags Sealur Moment -> materials
// @Security ApiKeyAuth
// @Description создание данных для материала
// @ModuleID createElasticity
// @Accept json
// @Produce json
// @Param elasticity body moment_model.CreateElasticityDTO true "elasticity info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/elasticity/ [post]
func (h *Handler) createElasticity(c *gin.Context) {
	var dto moment_model.CreateElasticityDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	elasticity := make([]*moment_api.Elasticity, 0, len(dto.Elasticity))
	for _, v := range dto.Elasticity {
		elasticity = append(elasticity, &moment_api.Elasticity{
			Temperature: v.Temperature,
			Elasticity:  v.Elasticity,
		})
	}

	_, err := h.materialsClient.CreateElasticity(c, &moment_api.CreateElasticityRequest{
		MarkId:     dto.MarkId,
		Elasticity: elasticity,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

// @Summary Update Elasticity
// @Tags Sealur Moment -> materials
// @Security ApiKeyAuth
// @Description обновление данных для материала
// @ModuleID updateElasticity
// @Accept json
// @Produce json
// @Param id path string true "elasticity id"
// @Param elasticity body moment_model.UpdateElasticityDTO true "elasticity info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/elasticity/{id} [put]
func (h *Handler) updateElasticity(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.UpdateElasticityDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.materialsClient.UpdateElasticity(c, &moment_api.UpdateElasticityRequest{
		Id:          id,
		MarkId:      dto.MarkId,
		Temperature: dto.Temperature,
		Elasticity:  dto.Elasticity,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

// @Summary Delete Elasticity
// @Tags Sealur Moment -> materials
// @Security ApiKeyAuth
// @Description Удаление данных для материала
// @ModuleID deleteElasticity
// @Accept json
// @Produce json
// @Param id path string true "elasticity id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/elasticity/{id} [delete]
func (h *Handler) deleteElasticity(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.materialsClient.DeleteElasticity(c, &moment_api.DeleteElasticityRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
