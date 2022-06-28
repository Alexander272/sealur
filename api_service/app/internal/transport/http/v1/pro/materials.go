package pro

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initMaterialsRoutes(api *gin.RouterGroup) {
	materials := api.Group("/materials")
	{
		materials.GET("/", h.getMaterials)
		materials = materials.Group("/", h.middleware.AccessForProAdmin)
		{
			materials.POST("/", h.createMaterials)
			materials.PUT("/:id", h.updateMaterials)
			materials.DELETE("/:id", h.deleteMaterials)
		}
	}
}

// @Summary Get Materials
// @Tags Sealur Pro -> materials
// @Security ApiKeyAuth
// @Description получение материалов (для опроса)
// @ModuleID getMaterials
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{data=[]proto.Materials}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/materials [get]
func (h *Handler) getMaterials(c *gin.Context) {
	mats, err := h.proClient.GetMaterials(c, &proto.GetMaterialsRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: mats.Materials})
}

// @Summary Create Materials
// @Tags Sealur Pro -> materials
// @Security ApiKeyAuth
// @Description создание материала (для опроса)
// @ModuleID createMaterials
// @Accept json
// @Produce json
// @Param data body models.MaterialsDTO true "materials info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/materials [post]
func (h *Handler) createMaterials(c *gin.Context) {
	var dto models.MaterialsDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	mat, err := h.proClient.CreateMaterials(c, &proto.CreateMaterialsRequest{
		Title:   dto.Title,
		TypeMat: dto.TypeMat,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/materials/%s", mat.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: mat.Id, Message: "Created"})
}

// @Summary Update Materials
// @Tags Sealur Pro -> materials
// @Security ApiKeyAuth
// @Description обновление материала (для опроса)
// @ModuleID updateMaterials
// @Accept json
// @Produce json
// @Param data body models.MaterialsDTO true "material info"
// @Param id path string true "material id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/materials/{id} [put]
func (h *Handler) updateMaterials(c *gin.Context) {
	var dto models.MaterialsDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	matId := c.Param("id")
	if matId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	mat, err := h.proClient.UpdateMaterials(c, &proto.UpdateMaterialsRequest{
		Id:      matId,
		Title:   dto.Title,
		TypeMat: dto.TypeMat,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: mat.Id, Message: "Updated"})
}

// @Summary Delete Materials
// @Tags Sealur Pro -> materials
// @Security ApiKeyAuth
// @Description удаление материала (для опроса)
// @ModuleID deleteMaterials
// @Accept json
// @Produce json
// @Param id path string true "material id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/materials/{id} [delete]
func (h *Handler) deleteMaterials(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	mat, err := h.proClient.DeleteMaterials(c, &proto.DeleteMaterialsRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: mat.Id, Message: "Deleted"})
}
