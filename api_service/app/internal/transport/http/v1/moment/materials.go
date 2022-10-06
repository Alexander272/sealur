package moment

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment/material_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initMaterialsRoutes(api *gin.RouterGroup) {
	materials := api.Group("/materials", h.middleware.UserIdentity)
	{
		materials.GET("/", h.getMaterials)
		materials = materials.Group("/", h.middleware.AccessForMomentAdmin)
		{
			materials.GET("/empty", h.getMaterialsWithIsEmpty)
			materials.GET("/:id", h.getMaterialsData)
			materials.POST("/", h.createMaterial)
			materials.PUT("/:id", h.updateMaterial)
			materials.DELETE("/:id", h.deleteMaterial)
		}
	}
}

// @Summary Get Materials
// @Tags Sealur Moment -> materials
// @Security ApiKeyAuth
// @Description получение типов материалов
// @ModuleID getMaterials
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{Data=[]material_api.Material}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/ [get]
func (h *Handler) getMaterials(c *gin.Context) {
	materials, err := h.materialsClient.GetMaterials(c, &material_api.GetMaterialsRequest{})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: materials.Materials, Count: len(materials.Materials)})
}

// @Summary Get Materials With Is Empty
// @Tags Sealur Moment -> materials
// @Security ApiKeyAuth
// @Description получение типов материалов с меткой о пустоту данных
// @ModuleID getMaterialsWithIsEmpty
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{Data=[]material_api.MaterialWithIsEmpty}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/empty [get]
func (h *Handler) getMaterialsWithIsEmpty(c *gin.Context) {
	materials, err := h.materialsClient.GetMaterialsWithIsEmpty(c, &material_api.GetMaterialsRequest{})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: materials.Materials, Count: len(materials.Materials)})
}

// @Summary Get Materials Data
// @Tags Sealur Moment -> materials
// @Security ApiKeyAuth
// @Description получение данных материала
// @ModuleID getMaterialsData
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{Data=material_api.MaterialsDataResponse}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/{id} [get]
func (h *Handler) getMaterialsData(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	materials, err := h.materialsClient.GetMaterialsData(c, &material_api.GetMaterialsDataRequest{MarkId: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: materials})
}

// @Summary Create Material
// @Tags Sealur Moment -> materials
// @Security ApiKeyAuth
// @Description создание материала
// @ModuleID createMaterial
// @Accept json
// @Produce json
// @Param material body moment_model.MaterialsDTO true "material info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/ [post]
func (h *Handler) createMaterial(c *gin.Context) {
	var dto moment_model.MaterialsDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	material, err := h.materialsClient.CreateMaterial(c, &material_api.CreateMaterialRequest{Title: dto.Title})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-moment/gasket/%s", material.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: material.Id, Message: "Created"})
}

// @Summary Update Material
// @Tags Sealur Moment -> materials
// @Security ApiKeyAuth
// @Description обновление материала
// @ModuleID updateMaterial
// @Accept json
// @Produce json
// @Param id path string true "material id"
// @Param material body moment_model.MaterialsDTO true "material info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/{id} [put]
func (h *Handler) updateMaterial(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.MaterialsDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.materialsClient.UpdateMaterial(c, &material_api.UpdateMaterialRequest{Id: id, Title: dto.Title})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Id: id, Message: "Updated"})
}

// @Summary Delete Material
// @Tags Sealur Moment -> materials
// @Security ApiKeyAuth
// @Description Удаление материала
// @ModuleID deleteMaterial
// @Accept json
// @Produce json
// @Param id path string true "material id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/material/{id} [delete]
func (h *Handler) deleteMaterial(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.materialsClient.DeleteMaterial(c, &material_api.DeleteMaterialRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Id: id, Message: "Deleted"})
}
