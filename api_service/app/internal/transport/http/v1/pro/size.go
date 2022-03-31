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
		sizes.DELETE("/all", h.deleteAllSize)
	}
}

// @Summary Get Sizes
// @Tags Sealur Pro -> sizes
// @Security ApiKeyAuth
// @Description получение размеров
// @ModuleID getSizes
// @Accept json
// @Produce json
// @Param flange query string true "flange"
// @Param typeFlId query string true "flange type id"
// @Param standId query string true "standarts id"
// @Param typePr query string true "type"
// @Success 200 {object} models.DataResponse{data=proto.SizeResponse}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes [get]
func (h *Handler) getSizes(c *gin.Context) {
	flange := c.Query("flange")
	if flange == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty flange", "empty flange param")
		return
	}
	typeFlId := c.Query("typeFlId")
	if typeFlId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty flange type id", "empty flange type id param")
		return
	}
	standId := c.Query("standId")
	if standId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty stand id", "empty standarts id param")
		return
	}
	typePr := c.Query("typePr")
	if typePr == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty type", "empty type lining")
		return
	}

	sizes, err := h.proClient.GetSizes(c, &proto.GetSizesRequest{
		Flange:   flange,
		TypeFlId: typeFlId,
		TypePr:   typePr,
		StandId:  standId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: sizes, Count: len(sizes.Sizes)})
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
		Flange:   dto.Flange,
		TypeFlId: dto.TypeFlId,
		Dn:       dto.Dn,
		Pn:       dto.Pn,
		TypePr:   dto.TypePr,
		StandId:  dto.StandId,
		D4:       dto.D4,
		D3:       dto.D3,
		D2:       dto.D2,
		D1:       dto.D1,
		H:        dto.H,
		S2:       dto.S2,
		S3:       dto.S3,
		Number:   dto.Number,
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
// @Success 200 {object} models.IdResponse
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
		Id:       id,
		Flange:   dto.Flange,
		TypeFlId: dto.TypeFlId,
		Dn:       dto.Dn,
		Pn:       dto.Pn,
		TypePr:   dto.TypePr,
		StandId:  dto.StandId,
		D4:       dto.D4,
		D3:       dto.D3,
		D2:       dto.D2,
		D1:       dto.D1,
		H:        dto.H,
		S2:       dto.S2,
		S3:       dto.S3,
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
// @Description удаление размеров
// @ModuleID deleteSize
// @Accept json
// @Produce json
// @Param id path string true "size id"
// @Param flange query string true "flange"
// @Success 200 {object} models.IdResponse
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

	flange := c.Query("flange")
	if flange == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty flange", "empty flange param")
		return
	}

	size, err := h.proClient.DeleteSize(c, &proto.DeleteSizeRequest{
		Id:     id,
		Flange: flange,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: size.Id, Message: "Deleted"})
}

// @Summary Delete Size
// @Tags Sealur Pro -> sizes
// @Security ApiKeyAuth
// @Description удаление всех размеров
// @ModuleID deleteSize
// @Accept json
// @Produce json
// @Param flange query string true "flange"
// @Param typePr query string true "type pr"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes/all [delete]
func (h *Handler) deleteAllSize(c *gin.Context) {
	flange := c.Query("flange")
	if flange == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty flange", "empty flange param")
		return
	}

	typePr := c.Query("typePr")
	if flange == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty type", "empty type param")
		return
	}

	size, err := h.proClient.DeleteAllSize(c, &proto.DeleteAllSizeRequest{
		Flange: flange,
		TypePr: typePr,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: size.Id, Message: "Deleted"})
}
