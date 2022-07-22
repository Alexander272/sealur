package moment

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto/moment_proto"
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
// @Success 200 {object} models.DataResponse{Data=[]moment_proto.Material}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/ [get]
func (h *Handler) getMaterials(c *gin.Context) {
	materials, err := h.materialsClient.GetMaterials(c, &moment_proto.GetMaterialsRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
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
// @Success 200 {object} models.DataResponse{Data=[]moment_proto.MaterialWithIsEmpty}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/empty [get]
func (h *Handler) getMaterialsWithIsEmpty(c *gin.Context) {
	materials, err := h.materialsClient.GetMaterialsWithIsEmpty(c, &moment_proto.GetMaterialsRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
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
// @Success 200 {object} models.DataResponse{Data=moment_proto.MaterialsDataResponse}
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

	materials, err := h.materialsClient.GetMaterialsData(c, &moment_proto.GetMaterialsDataRequest{MarkId: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
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
// @Param material body models.MomentMaterialsDTO true "material info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/materials/ [post]
func (h *Handler) createMaterial(c *gin.Context) {
	var dto models.MomentMaterialsDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	material, err := h.materialsClient.CreateMaterial(c, &moment_proto.CreateMaterialRequest{Title: dto.Title})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
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
// @Param material body models.MomentMaterialsDTO true "material info"
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

	var dto models.MomentMaterialsDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.materialsClient.UpdateMaterial(c, &moment_proto.UpdateMaterialRequest{Id: id, Title: dto.Title})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
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

	_, err := h.materialsClient.DeleteMaterial(c, &moment_proto.DeleteMaterialRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Id: id, Message: "Deleted"})
}
