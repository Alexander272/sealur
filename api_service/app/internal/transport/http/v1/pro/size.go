package pro

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initSizeRoutes(api *gin.RouterGroup) {
	sizes := api.Group("/sizes")
	{
		sizes.GET("/", h.getSizes)
		sizes.POST("/", h.createSize)
		sizes.PUT("/:id", h.updateSize)
		sizes.DELETE("/:id", h.deleteSize)
	}
}

// @Summary Get Sizes
// @Tags Sealur Pro -> sizes
// @Security ApiKeyAuth
// @Description получение размеров
// @ModuleID getSizes
// @Accept json
// @Produce json
// @Param data body models.GetSizesDTO true "info for size"
// @Success 200 {object} models.DataResponse{data=[]proto.Size}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes [get]
func (h *Handler) getSizes(c *gin.Context) {
	var dto models.GetSizesDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	sizes, err := h.proClient.GetSizes(c, &proto.GetSizesRequest{
		Flange:  dto.Flange,
		TypeFl:  dto.TypeFl,
		TypePr:  dto.TypePr,
		StandId: dto.StandId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: sizes.Sizes})
}

// @Summary Create Size
// @Tags Sealur Pro -> sizes
// @Security ApiKeyAuth
// @Description создание размеров
// @ModuleID createSize
// @Accept json
// @Produce json
// @Param data body models.SizesDTO true "size info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes [post]
func (h *Handler) createSize(c *gin.Context) {
	var dto models.SizesDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	size, err := h.proClient.CreateSize(c, &proto.CreateSizeRequest{
		Flange:  dto.Flange,
		TypeFl:  dto.TypeFl,
		Dn:      dto.Dn,
		Pn:      dto.Pn,
		TypePr:  dto.TypePr,
		StandId: dto.StandId,
		D4:      dto.D4,
		D3:      dto.D3,
		D2:      dto.D2,
		D1:      dto.D1,
		H:       dto.H,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/sizes/%s", size.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: size.Id, Message: "Created"})
}

// @Summary Update Size
// @Tags Sealur Pro -> sizes
// @Security ApiKeyAuth
// @Description обновление размеров
// @ModuleID updateSize
// @Accept json
// @Produce json
// @Param data body models.SizesDTO true "size info"
// @Param id path string true "size id"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes/{id} [put]
func (h *Handler) updateSize(c *gin.Context) {
	var dto models.SizesDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	size, err := h.proClient.UpdateSize(c, &proto.UpdateSizeRequest{
		Id:      id,
		Flange:  dto.Flange,
		TypeFl:  dto.TypeFl,
		Dn:      dto.Dn,
		Pn:      dto.Pn,
		TypePr:  dto.TypePr,
		StandId: dto.StandId,
		D4:      dto.D4,
		D3:      dto.D3,
		D2:      dto.D2,
		D1:      dto.D1,
		H:       dto.H,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: size.Id, Message: "Updated"})
}

// @Summary Delete Size
// @Tags Sealur Pro -> sizes
// @Security ApiKeyAuth
// @Description обновление размеров
// @ModuleID deleteSize
// @Accept json
// @Produce json
// @Param data body models.DeleteSizeDTO true "info for size"
// @Param id path string true "size id"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes/{id} [delete]
func (h *Handler) deleteSize(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto models.DeleteSizeDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	size, err := h.proClient.DeleteSize(c, &proto.DeleteSizeRequest{
		Id:     id,
		Flange: dto.Flange,
		TypeFl: dto.TypeFl,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: size.Id, Message: "Deleted"})
}
