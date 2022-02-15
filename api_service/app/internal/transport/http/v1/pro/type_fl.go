package pro

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initTypeFlRoutes(api *gin.RouterGroup) {
	fl := api.Group("/type-fl")
	{
		fl.GET("/", h.getTypeFl)
		fl.GET("/all", h.getAllTypeFl)
		fl.POST("/", h.createTypeFl)
		fl.PUT("/:id", h.updateTypeFl)
		fl.DELETE("/:id", h.deleteTypeFl)
	}
}

// @Summary Get Type Flange
// @Tags Sealur Pro -> type-fl
// @Security ApiKeyAuth
// @Description получение базовых типов фланца
// @ModuleID getTypeFl
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{data=[]proto.TypeFl}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/type-fl [get]
func (h *Handler) getTypeFl(c *gin.Context) {
	fl, err := h.proClient.GetTypeFl(c, &proto.GetTypeFlRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wront")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: fl.TypeFl})
}

// @Summary Get All Type Flange
// @Tags Sealur Pro -> type-fl
// @Security ApiKeyAuth
// @Description получение всех типов фланца
// @ModuleID getAllTypeFl
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{data=[]proto.TypeFl}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/type-fl/all [get]
func (h *Handler) getAllTypeFl(c *gin.Context) {
	fl, err := h.proClient.GetAllTypeFl(c, &proto.GetTypeFlRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: fl.TypeFl})
}

// @Summary Create Type Flange
// @Tags Sealur Pro -> type-fl
// @Security ApiKeyAuth
// @Description создание типа фланца
// @ModuleID createTypeFl
// @Accept json
// @Produce json
// @Param data body models.TypeFlDTO true "type flange info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/type-fl [post]
func (h *Handler) createTypeFl(c *gin.Context) {
	var dto models.TypeFlDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	request := &proto.CreateTypeFlRequest{
		Title: dto.Title,
		Descr: dto.Desc,
		Short: dto.Short,
		Basis: dto.Basis,
	}

	fl, err := h.proClient.CreateTypeFl(c, request)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/type-fl/%s", fl.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: fl.Id, Message: "Created"})
}

// @Summary Update Type Flange
// @Tags Sealur Pro -> type-fl
// @Security ApiKeyAuth
// @Description обновление типа фланца
// @ModuleID updateTypeFl
// @Accept json
// @Produce json
// @Param id path string true "type flange id"
// @Param data body models.TypeFlDTO true "type flange info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/type-fl/{id} [put]
func (h *Handler) updateTypeFl(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto models.TypeFlDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	request := &proto.UpdateTypeFlRequest{
		Id:    id,
		Title: dto.Title,
		Descr: dto.Desc,
		Short: dto.Short,
		Basis: dto.Basis,
	}

	fl, err := h.proClient.UpdateTypeFl(c, request)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: fl.Id, Message: "Updated"})
}

// @Summary Delete Type Flange
// @Tags Sealur Pro -> type-fl
// @Security ApiKeyAuth
// @Description удаление типа фланца
// @ModuleID deleteTypeFl
// @Accept json
// @Produce json
// @Param id path string true "type flange id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/type-fl/{id} [delete]
func (h *Handler) deleteTypeFl(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	fl, err := h.proClient.DeleteTypeFl(c, &proto.DeleteTypeFlRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: fl.Id, Message: "Deleted"})
}
