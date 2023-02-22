package pro

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/pro_model"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initTypeFlRoutes(api *gin.RouterGroup) {
	fl := api.Group("/flange-types")
	{
		fl.GET("/", h.getTypeFl)
		fl.GET("/all", h.getAllTypeFl)
		fl = fl.Group("/", h.middleware.UserIdentity, h.middleware.AccessForProAdmin)
		{
			fl.POST("/", h.createTypeFl)
			fl.PUT("/:id", h.updateTypeFl)
			fl.DELETE("/:id", h.deleteTypeFl)
		}
	}
}

// @Summary Get Type Flange
// @Tags Sealur Pro -> flange-types
// @Description получение базовых типов фланца
// @ModuleID getTypeFl
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{data=[]pro_api.TypeFl}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/flange-types [get]
func (h *Handler) getTypeFl(c *gin.Context) {
	fl, err := h.proClient.GetTypeFl(c, &pro_api.GetTypeFlRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wront")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: fl.TypeFl})
}

// @Summary Get All Type Flange
// @Tags Sealur Pro -> flange-types
// @Description получение всех типов фланца
// @ModuleID getAllTypeFl
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{data=[]pro_api.TypeFl}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/flange-types/all [get]
func (h *Handler) getAllTypeFl(c *gin.Context) {
	fl, err := h.proClient.GetAllTypeFl(c, &pro_api.GetTypeFlRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: fl.TypeFl})
}

// @Summary Create Type Flange
// @Tags Sealur Pro -> flange-types
// @Security ApiKeyAuth
// @Description создание типа фланца
// @ModuleID createTypeFl
// @Accept json
// @Produce json
// @Param data body pro_model.TypeFlDTO true "type flange info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/flange-types [post]
func (h *Handler) createTypeFl(c *gin.Context) {
	var dto pro_model.TypeFlDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	request := &pro_api.CreateTypeFlRequest{
		Title: dto.Title,
		Descr: dto.Descr,
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
// @Tags Sealur Pro -> flange-types
// @Security ApiKeyAuth
// @Description обновление типа фланца
// @ModuleID updateTypeFl
// @Accept json
// @Produce json
// @Param id path string true "type flange id"
// @Param data body pro_model.TypeFlDTO true "type flange info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/flange-types/{id} [put]
func (h *Handler) updateTypeFl(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto pro_model.TypeFlDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	request := &pro_api.UpdateTypeFlRequest{
		Id:    id,
		Title: dto.Title,
		Descr: dto.Descr,
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
// @Tags Sealur Pro -> flange-types
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
// @Router /sealur-pro/flange-types/{id} [delete]
func (h *Handler) deleteTypeFl(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	fl, err := h.proClient.DeleteTypeFl(c, &pro_api.DeleteTypeFlRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: fl.Id, Message: "Deleted"})
}
