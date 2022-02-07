package pro

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initFlangeRoutes(api *gin.RouterGroup) {
	stand := api.Group("/flanges")
	{
		stand.GET("/", h.GetFlanges)
		stand.POST("/", h.notImplemented)
		stand.PUT("/:id", h.notImplemented)
		stand.DELETE("/:id", h.notImplemented)
	}
}

// @Summary Get Flanges
// @Tags Sealur Pro -> flanges
// @Security ApiKeyAuth
// @Description получение всех стандартов на фланцы
// @ModuleID getFlanges
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{data=[]proto.Flange}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/flanges [get]
func (h *Handler) GetFlanges(c *gin.Context) {
	fl, err := h.proClient.GetAllFlanges(c, &proto.GetAllFlangeRequest{})
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
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/flanges [post]
func (h *Handler) CreateFlange(c *gin.Context) {
	var dto models.FlangeDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id, err := h.proClient.CreateFlange(c, &proto.CreateFlangeRequest{Title: dto.Title, Short: dto.Short})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/flanges/%s", id.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: id.Id, Message: "Created"})
}

// @Summary Update Flange
// @Tags Sealur Pro -> flanges
// @Security ApiKeyAuth
// @Description обновление стандарта на фланцы
// @ModuleID updateFlange
// @Accept json
// @Produce json
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/flanges/{id} [put]
func (h *Handler) UpdateFlange(c *gin.Context) {
	var dto models.FlangeDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	flId := c.Param("id")
	if flId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	fl, err := h.proClient.UpdateFlange(c, &proto.UpdateFlangeRequest{Title: dto.Title, Short: dto.Short})
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
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/flanges/{id} [delete]
func (h *Handler) DeleteFlange(c *gin.Context) {
	flId := c.Param("id")
	if flId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	fl, err := h.proClient.DeleteFlange(c, &proto.DeleteFlangeRequest{Id: flId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: fl.Id, Message: "Deleted"})
}
