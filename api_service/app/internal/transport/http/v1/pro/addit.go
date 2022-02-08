package pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initAdditRoutes(api *gin.RouterGroup) {
	addit := api.Group("/additionals")
	{
		addit.GET("/", h.getAddit)
		addit.POST("/", h.createAddit)
		addit.PATCH("/:id/mat", h.updateMat)
		addit.PATCH("/:id/mod", h.updateMod)
		addit.PATCH("/:id/temp", h.updateTemp)
		addit.PATCH("/:id/moun", h.updateMoun)
		addit.PATCH("/:id/grap", h.updateGrap)
		addit.PATCH("/:id/fl", h.updateTypeFl)
	}
}

// @Summary Get Addit
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description получение всех дополнительных сведенний (список материалов, креплений и тд)
// @ModuleID getAddit
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{data=[]proto.Additional}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals [get]
func (h *Handler) getAddit(c *gin.Context) {
	add, err := h.proClient.GetAllAdditional(c, &proto.GetAllAddRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: add.Additionals})
}

// @Summary Create Addit
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description создание дополнительных сведенний
// @ModuleID createAddit
// @Accept json
// @Produce json
// @Param data body models.CreateAdditDTO true "addit info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals [post]
func (h *Handler) createAddit(c *gin.Context) {
	var dto models.CreateAdditDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.proClient.CreateAdditional(c, &proto.CreateAddRequest{
		Materials:   dto.Materials,
		Mod:         dto.Mod,
		Temperature: dto.Temperature,
		Mounting:    dto.Mounting,
		Graphite:    dto.Graphite,
		TypeFl:      dto.TypeFl,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

// @Summary Update Mat
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление материалов
// @ModuleID updateMat
// @Accept json
// @Produce json
// @Param data body models.UpdateMatDTO true "additional materials info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/mat [patch]
func (h *Handler) updateMat(c *gin.Context) {
	var dto models.UpdateMatDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateMat(c, &proto.UpdateAddMatRequest{Id: id, Materials: dto.Materials})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated materials"})
}

// @Summary Update Mod
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление модифицирующего элемента
// @ModuleID updateMod
// @Accept json
// @Produce json
// @Param data body models.UpdateModDTO true "additional modification info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/mod [patch]
func (h *Handler) updateMod(c *gin.Context) {
	var dto models.UpdateModDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateMod(c, &proto.UpdateAddModRequest{Id: id, Mod: dto.Mod})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated modifications"})
}

// @Summary Update Temp
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление температур
// @ModuleID updateTemp
// @Accept json
// @Produce json
// @Param data body models.UpdateTempDTO true "additional temperature info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/temp [patch]
func (h *Handler) updateTemp(c *gin.Context) {
	var dto models.UpdateTempDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateTemp(c, &proto.UpdateAddTemRequest{Id: id, Temperature: dto.Temperature})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated temperature"})
}

// @Summary Update Moun
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление крепления на вертикальном фланце
// @ModuleID updateMoun
// @Accept json
// @Produce json
// @Param data body models.UpdateMounDTO true "additional mounting info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/moun [patch]
func (h *Handler) updateMoun(c *gin.Context) {
	var dto models.UpdateMounDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateMoun(c, &proto.UpdateAddMounRequest{Id: id, Mounting: dto.Mounting})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated mounting"})
}

// @Summary Update Grap
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление графита
// @ModuleID updateGrap
// @Accept json
// @Produce json
// @Param data body models.UpdateGrapDTO true "additional graphite info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/grap [patch]
func (h *Handler) updateGrap(c *gin.Context) {
	var dto models.UpdateGrapDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateGrap(c, &proto.UpdateAddGrapRequest{Id: id, Graphite: dto.Graphite})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated graphite"})
}

// @Summary Update Type Flange
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление типа фланца
// @ModuleID updateTypeFl
// @Accept json
// @Produce json
// @Param data body models.UpdateFlDTO true "additional flange info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/fl [patch]
func (h *Handler) updateTypeFl(c *gin.Context) {
	var dto models.UpdateFlDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateTypeFl(c, &proto.UpdateAddTypeFlRequest{Id: id, TypeFl: dto.TypeFl})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Update type flanges"})
}
