package pro

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initPutgmRoutes(api *gin.RouterGroup) {
	putgm := api.Group("/putgm")
	{
		putgm.GET("/", h.getPutgm)
		putgm = putgm.Group("/", h.middleware.AccessForProAdmin)
		{
			putgm.POST("/", h.createPutgm)
			putgm.PUT("/:id", h.updatePutgm)
			putgm.DELETE("/:id", h.deletePutgm)
		}
	}
}

// @Summary Get Putgm
// @Tags Sealur Pro -> putgm
// @Description получение прокладок путгм
// @ModuleID getPutgm
// @Accept json
// @Produce json
// @Param form query string true "form"
// @Param flangeId query string true "flange id"
// @Success 200 {object} models.DataResponse{data=[]proto.Putgm}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putgm [get]
func (h *Handler) getPutgm(c *gin.Context) {
	form := c.Query("form")
	if form == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty form param")
		return
	}
	flangeId := c.Query("flangeId")
	if flangeId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty flangeId param")
		return
	}

	putgm, err := h.proClient.GetPutgm(c, &proto.GetPutgmRequest{
		Form:     form,
		FlangeId: flangeId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: putgm.Putgm, Count: len(putgm.Putgm)})
}

// @Summary Create Putgm
// @Tags Sealur Pro -> putgm
// @Security ApiKeyAuth
// @Description создание прокладки путгм
// @ModuleID createPutgm
// @Accept json
// @Produce json
// @Param data body models.PutgmDTO true "putgm info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putgm [post]
func (h *Handler) createPutgm(c *gin.Context) {
	var dto models.PutgmDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	putg, err := h.proClient.CreatePutgm(c, &proto.CreatePutgmRequest{
		FlangeId:     dto.FlangeId,
		TypeFlId:     dto.TypeFlId,
		TypePr:       dto.TypePr,
		Form:         dto.Form,
		Construction: dto.Construction,
		Temperatures: dto.Temperatures,
		Basis:        dto.Basis,
		Obturator:    dto.Obturator,
		Coating:      dto.Coating,
		Mounting:     dto.Mounting,
		Graphite:     dto.Graphite,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/putgm/%s", putg.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: putg.Id, Message: "Created"})
}

// @Summary Update Putgm
// @Tags Sealur Pro -> putgm
// @Security ApiKeyAuth
// @Description обновление прокладки путгм
// @ModuleID updatePutgm
// @Accept json
// @Produce json
// @Param data body models.PutgmDTO true "putgm info"
// @Param id path string true "putgm id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putgm/{id} [put]
func (h *Handler) updatePutgm(c *gin.Context) {
	var dto models.PutgmDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	putg, err := h.proClient.UpdatePutgm(c, &proto.UpdatePutgmRequest{
		Id:           id,
		FlangeId:     dto.FlangeId,
		TypeFlId:     dto.TypeFlId,
		TypePr:       dto.TypePr,
		Form:         dto.Form,
		Construction: dto.Construction,
		Temperatures: dto.Temperatures,
		Basis:        dto.Basis,
		Obturator:    dto.Obturator,
		Coating:      dto.Coating,
		Mounting:     dto.Mounting,
		Graphite:     dto.Graphite,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: putg.Id, Message: "Updated"})
}

// @Summary Delete Putgm
// @Tags Sealur Pro -> putgm
// @Security ApiKeyAuth
// @Description удаление прокладки путгм
// @ModuleID deletePutgm
// @Accept json
// @Produce json
// @Param id path string true "putgm id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putgm/{id} [delete]
func (h *Handler) deletePutgm(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	putg, err := h.proClient.DeletePutgm(c, &proto.DeletePutgmRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: putg.Id, Message: "Deleted"})
}
