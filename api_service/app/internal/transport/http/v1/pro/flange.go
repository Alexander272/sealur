package pro

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/pro_model"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initFlangeRoutes(api *gin.RouterGroup) {
	flanges := api.Group("/flanges")
	{
		flanges.GET("/", h.getFlanges)
		flanges = flanges.Group("/", h.middleware.UserIdentity, h.middleware.AccessForProAdmin)
		{
			flanges.POST("/", h.createFlange)
			flanges.PUT("/:id", h.updateFlange)
			flanges.DELETE("/:id", h.deleteFlange)
		}
	}
}

// @Summary Get Flanges
// @Tags Sealur Pro -> flanges
// @Description получение всех стандартов на фланцы
// @ModuleID getFlanges
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{data=[]pro_api.Flange}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/flanges [get]
func (h *Handler) getFlanges(c *gin.Context) {
	fl, err := h.proClient.GetAllFlanges(c, &pro_api.GetAllFlangeRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: fl.Flanges})
}

// @Summary Create Flange
// @Tags Sealur Pro -> flanges
// @Security ApiKeyAuth
// @Description создание стандарта на фланцы
// @ModuleID createFlange
// @Accept json
// @Produce json
// @Param data body pro_model.FlangeDTO true "flange info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/flanges [post]
func (h *Handler) createFlange(c *gin.Context) {
	var dto pro_model.FlangeDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	fl, err := h.proClient.CreateFlange(c, &pro_api.CreateFlangeRequest{Title: dto.Title, Short: dto.Short})
	if err != nil {
		if errors.Is(err, models.ErrFlangeAlreadyExists) {
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
			return
		}
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/flanges/%s", fl.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: fl.Id, Message: "Created"})
}

// @Summary Update Flange
// @Tags Sealur Pro -> flanges
// @Security ApiKeyAuth
// @Description обновление стандарта на фланцы
// @ModuleID updateFlange
// @Accept json
// @Produce json
// @Param data body pro_model.FlangeDTO true "flange info"
// @Param id path string true "flange id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/flanges/{id} [put]
func (h *Handler) updateFlange(c *gin.Context) {
	var dto pro_model.FlangeDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	flId := c.Param("id")
	if flId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	fl, err := h.proClient.UpdateFlange(c, &pro_api.UpdateFlangeRequest{Title: dto.Title, Short: dto.Short})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: fl.Id, Message: "Updated"})
}

// @Summary Delete Flange
// @Tags Sealur Pro -> flanges
// @Security ApiKeyAuth
// @Description удаление стандарта на фланцы
// @ModuleID deleteFlange
// @Accept json
// @Produce json
// @Param id path string true "flange id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/flanges/{id} [delete]
func (h *Handler) deleteFlange(c *gin.Context) {
	flId := c.Param("id")
	if flId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	fl, err := h.proClient.DeleteFlange(c, &pro_api.DeleteFlangeRequest{Id: flId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: fl.Id, Message: "Deleted"})
}
