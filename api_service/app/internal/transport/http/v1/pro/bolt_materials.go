package pro

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initBoltMaterialsRoutes(api *gin.RouterGroup) {
	materials := api.Group("/bolt-materials")
	{
		materials.GET("/", h.getBoltMaterials)
		materials = materials.Group("/", h.middleware.AccessForProAdmin)
		{
			materials.POST("/", h.createBoltMaterials)
			materials.PUT("/:id", h.updateBoltMaterials)
			materials.DELETE("/:id", h.deleteBoltMaterials)
		}
	}
}

// @Summary Get Bolt Materials
// @Tags Sealur Pro -> bolt materials
// @Description получение материалов болтов/шпилек (для опроса)
// @ModuleID getBoltMaterials
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{data=[]proto.BoltMaterials}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/bolt-materials [get]
func (h *Handler) getBoltMaterials(c *gin.Context) {
	mats, err := h.proClient.GetBoltMaterials(c, &proto.GetBoltMaterialsRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: mats.Materials})
}

// @Summary Create Bolt Materials
// @Tags Sealur Pro -> bolt materials
// @Security ApiKeyAuth
// @Description создание материала болтов/шпилек (для опроса)
// @ModuleID createBoltMaterials
// @Accept json
// @Produce json
// @Param data body models.BoltMaterialsDTO true "bolt materials info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/bolt-materials [post]
func (h *Handler) createBoltMaterials(c *gin.Context) {
	var dto models.BoltMaterialsDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	mat, err := h.proClient.CreateBoltMaterials(c, &proto.CreateBoltMaterialsRequest{
		Title:    dto.Title,
		FlangeId: dto.FlangeId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/bolt-materials/%s", mat.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: mat.Id, Message: "Created"})
}

// @Summary Update Bolt Materials
// @Tags Sealur Pro -> bolt materials
// @Security ApiKeyAuth
// @Description обновление материала болтов/шпилек (для опроса)
// @ModuleID updateBoltMaterials
// @Accept json
// @Produce json
// @Param data body models.BoltMaterialsDTO true "bolt material info"
// @Param id path string true "bolt material id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/bolt-materials/{id} [put]
func (h *Handler) updateBoltMaterials(c *gin.Context) {
	var dto models.BoltMaterialsDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	matId := c.Param("id")
	if matId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	mat, err := h.proClient.UpdateBoltMaterials(c, &proto.UpdateBoltMaterialsRequest{
		Id:       matId,
		Title:    dto.Title,
		FlangeId: dto.FlangeId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: mat.Id, Message: "Updated"})
}

// @Summary Delete Bolt Materials
// @Tags Sealur Pro -> bolt materials
// @Security ApiKeyAuth
// @Description удаление материала болтов/шпилек (для опроса)
// @ModuleID deleteBoltMaterials
// @Accept json
// @Produce json
// @Param id path string true "bolt material id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/bolt-materials/{id} [delete]
func (h *Handler) deleteBoltMaterials(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	mat, err := h.proClient.DeleteBoltMaterials(c, &proto.DeleteBoltMaterialsRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: mat.Id, Message: "Deleted"})
}
