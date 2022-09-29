package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initStandartsRoutes(api *gin.RouterGroup) {
	standarts := api.Group("/standarts", h.middleware.UserIdentity)
	{
		standarts.GET("/", h.getStandarts)
		standarts.GET("/size", h.getStandartsWithSize)
		standarts = standarts.Group("/", h.middleware.AccessForMomentAdmin)
		{
			standarts.POST("/", h.createStandart)
			standarts.PUT("/:id", h.updateStandart)
			standarts.DELETE("/:id", h.deleteStandart)
		}
	}
}

// @Summary Get Standarts
// @Tags Sealur Moment -> standarts
// @Security ApiKeyAuth
// @Description получение стандартов
// @ModuleID getStandarts
// @Accept json
// @Produce json
// @Param typeId query string true "type id"
// @Success 200 {object} models.DataResponse{Data=[]moment_api.Standart}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/standarts/ [get]
func (h *Handler) getStandarts(c *gin.Context) {
	typeId := c.Query("typeId")
	if typeId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty typeId param")
		return
	}

	standarts, err := h.flangeClient.GetStandarts(c, &moment_api.GetStandartsRequest{TypeId: typeId})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: standarts.Standarts, Count: len(standarts.Standarts)})
}

// @Summary Get Standarts With Size
// @Tags Sealur Moment -> standarts
// @Security ApiKeyAuth
// @Description получение стандартов
// @ModuleID getStandartsWithSize
// @Accept json
// @Produce json
// @Param typeId query string true "type id"
// @Success 200 {object} models.DataResponse{Data=[]moment_api.StandartWithSize}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/standarts/size [get]
func (h *Handler) getStandartsWithSize(c *gin.Context) {
	typeId := c.Query("typeId")
	if typeId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty typeId param")
		return
	}

	standarts, err := h.flangeClient.GetStandartsWithSize(c, &moment_api.GetStandartsRequest{TypeId: typeId})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: standarts.Standarts, Count: len(standarts.Standarts)})
}

// @Summary Create Standart
// @Tags Sealur Moment -> standarts
// @Security ApiKeyAuth
// @Description создание болта
// @ModuleID createStandart
// @Accept json
// @Produce json
// @Param standart body moment_model.MomentStandartDTO true "standart info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/standarts/ [post]
func (h *Handler) createStandart(c *gin.Context) {
	var dto moment_model.StandartDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	stand, err := h.flangeClient.CreateStandart(c, &moment_api.CreateStandartRequest{
		Title:          dto.Title,
		TypeId:         dto.TypeId,
		TitleDn:        dto.TitleDn,
		TitlePn:        dto.TitlePn,
		IsNeedRow:      dto.IsNeedRow,
		Rows:           dto.Rows,
		HasDesignation: dto.HasDesignation,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Id: stand.Id, Message: "Created"})
}

// @Summary Update Standart
// @Tags Sealur Moment -> standarts
// @Security ApiKeyAuth
// @Description обновление среды
// @ModuleID updateStandart
// @Accept json
// @Produce json
// @Param id path string true "standart id"
// @Param standart body moment_model.StandartDTO true "standart info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/standarts/{id} [put]
func (h *Handler) updateStandart(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.StandartDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.flangeClient.UpdateStandart(c, &moment_api.UpdateStandartRequest{
		Id:             id,
		Title:          dto.Title,
		TypeId:         dto.TypeId,
		TitleDn:        dto.TitleDn,
		TitlePn:        dto.TitlePn,
		IsNeedRow:      dto.IsNeedRow,
		Rows:           dto.Rows,
		IsInch:         dto.IsInch,
		HasDesignation: dto.HasDesignation,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

// @Summary Delete Standart
// @Tags Sealur Moment -> standarts
// @Security ApiKeyAuth
// @Description Удаление болта
// @ModuleID deleteStandart
// @Accept json
// @Produce json
// @Param id path string true "standart id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/standarts/{id} [delete]
func (h *Handler) deleteStandart(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.flangeClient.DeleteStandart(c, &moment_api.DeleteStandartRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
