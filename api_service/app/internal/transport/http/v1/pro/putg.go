package pro

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initPutgRoutes(api *gin.RouterGroup) {
	putg := api.Group("/putg")
	{
		putg.GET("/", h.getPutg)
		putg = putg.Group("/", h.middleware.AccessForProAdmin)
		{
			putg.POST("/", h.createPutg)
			putg.PUT("/:id", h.updatePutg)
			putg.DELETE("/:id", h.deletePutg)
		}
	}
}

// @Summary Get Putg
// @Tags Sealur Pro -> putg
// @Security ApiKeyAuth
// @Description получение прокладок путг
// @ModuleID getPutg
// @Accept json
// @Produce json
// @Param form query string true "form"
// @Param flangeId query string true "flange id"
// @Success 200 {object} models.DataResponse{data=[]proto.Putg}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putg [get]
func (h *Handler) getPutg(c *gin.Context) {
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

	putg, err := h.proClient.GetPutg(c, &proto.GetPutgRequest{
		Form:     form,
		FlangeId: flangeId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: putg.Putg, Count: len(putg.Putg)})
}

// @Summary Create Putg
// @Tags Sealur Pro -> putg
// @Security ApiKeyAuth
// @Description создание прокладки путг
// @ModuleID createPutg
// @Accept json
// @Produce json
// @Param data body models.PutgDTO true "putg info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putg [post]
func (h *Handler) createPutg(c *gin.Context) {
	var dto models.PutgDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	putg, err := h.proClient.CreatePutg(c, &proto.CreatePutgRequest{
		FlangeId:     dto.FlangeId,
		TypeFlId:     dto.TypeFlId,
		TypePr:       dto.TypePr,
		Form:         dto.Form,
		Construction: dto.Construction,
		Temperatures: dto.Temperatures,
		Reinforce:    dto.Reinforce,
		Obturator:    dto.Obturator,
		ILimiter:     dto.ILimiter,
		OLimiter:     dto.OLimiter,
		Coating:      dto.Coating,
		Mounting:     dto.Mounting,
		Graphite:     dto.Graphite,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/putg/%s", putg.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: putg.Id, Message: "Created"})
}

// @Summary Update Putg
// @Tags Sealur Pro -> putg
// @Security ApiKeyAuth
// @Description обновление прокладки путг
// @ModuleID updatePutg
// @Accept json
// @Produce json
// @Param data body models.PutgDTO true "putg info"
// @Param id path string true "putg id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putg/{id} [put]
func (h *Handler) updatePutg(c *gin.Context) {
	var dto models.PutgDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	putg, err := h.proClient.UpdatePutg(c, &proto.UpdatePutgRequest{
		Id:           id,
		FlangeId:     dto.FlangeId,
		TypeFlId:     dto.TypeFlId,
		TypePr:       dto.TypePr,
		Form:         dto.Form,
		Construction: dto.Construction,
		Temperatures: dto.Temperatures,
		Reinforce:    dto.Reinforce,
		Obturator:    dto.Obturator,
		ILimiter:     dto.ILimiter,
		OLimiter:     dto.OLimiter,
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

// @Summary Delete Putg
// @Tags Sealur Pro -> putg
// @Security ApiKeyAuth
// @Description удаление прокладки путг
// @ModuleID deletePutg
// @Accept json
// @Produce json
// @Param id path string true "putg id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putg/{id} [delete]
func (h *Handler) deletePutg(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	putg, err := h.proClient.DeletePutg(c, &proto.DeletePutgRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: putg.Id, Message: "Deleted"})
}
