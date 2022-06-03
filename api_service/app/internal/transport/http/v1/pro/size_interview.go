package pro

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initSizeIntRoutes(api *gin.RouterGroup) {
	size := api.Group("/size-interview")
	{
		size.GET("/", h.getSizesInt)
		size.POST("/", h.createSizeInt)
		size.PUT("/:id", h.updateSizeInt)
		size.DELETE("/:id", h.deleteSizeInt)
	}
}

// @Summary Get Sizes Int
// @Tags Sealur Pro -> sizes-interview
// @Security ApiKeyAuth
// @Description получение размеров (для опроса)
// @ModuleID getSizesInt
// @Accept json
// @Produce json
// @Param flange query string true "flange"
// @Param typeFlId query string true "flange type id"
// @Success 200 {object} models.DataResponse{data=proto.SizeIntResponse}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes-interview [get]
func (h *Handler) getSizesInt(c *gin.Context) {
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

	sizes, err := h.proClient.GetSizeInt(c, &proto.GetSizesIntRequest{
		FlangeId: flange,
		TypeFl:   typeFlId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: sizes, Count: len(sizes.Sizes)})
}

// @Summary Create Size Int
// @Tags Sealur Pro -> sizes-interview
// @Security ApiKeyAuth
// @Description создание размеров (для опроса)
// @ModuleID createSizeInt
// @Accept json
// @Produce json
// @Param data body models.SizeIntDTO true "size int info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes-interview [post]
func (h *Handler) createSizeInt(c *gin.Context) {
	var dto models.SizeIntDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	size, err := h.proClient.CreateSizeInt(c, &proto.CreateSizeIntRequest{
		FlangeId:  dto.Flange,
		TypeFl:    dto.TypeFlId,
		Dy:        dto.Dy,
		Py:        dto.Py,
		D1:        dto.D1,
		D2:        dto.D2,
		D:         dto.D,
		H1:        dto.H1,
		H2:        dto.H2,
		BoltId:    dto.BoltId,
		CountBolt: dto.CountBolt,
		Number:    dto.Count,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/sizes-interview/%s", size.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: size.Id, Message: "Created"})
}

// @Summary Update Size Int
// @Tags Sealur Pro -> sizes-interview
// @Security ApiKeyAuth
// @Description обновление размеров (для опроса)
// @ModuleID updateSizeInt
// @Accept json
// @Produce json
// @Param data body models.SizeIntDTO true "size int info"
// @Param id path string true "size int id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes-interview/{id} [put]
func (h *Handler) updateSizeInt(c *gin.Context) {
	var dto models.SizeIntDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	size, err := h.proClient.UpdateSizeInt(c, &proto.UpdateSizeIntRequest{
		Id:        id,
		FlangeId:  dto.Flange,
		TypeFl:    dto.TypeFlId,
		Dy:        dto.Dy,
		Py:        dto.Py,
		D1:        dto.D1,
		D2:        dto.D2,
		D:         dto.D,
		H1:        dto.H1,
		H2:        dto.H2,
		BoltId:    dto.BoltId,
		CountBolt: dto.CountBolt,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: size.Id, Message: "Updated"})
}

// @Summary Delete Size Int
// @Tags Sealur Pro -> sizes-interview
// @Security ApiKeyAuth
// @Description удаление размеров (для опроса)
// @ModuleID deleteSizeInt
// @Accept json
// @Produce json
// @Param id path string true "size id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes-interview/{id} [delete]
func (h *Handler) deleteSizeInt(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	size, err := h.proClient.DeleteSizeInt(c, &proto.DeleteSizeIntRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: size.Id, Message: "Deleted"})
}
