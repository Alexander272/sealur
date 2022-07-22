package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto/moment_proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initTypeFlangeRoutes(api *gin.RouterGroup) {
	typeFlnage := api.Group("/type-flange", h.middleware.UserIdentity)
	{
		typeFlnage.GET("/", h.getTypeFlange)
		typeFlnage = typeFlnage.Group("/", h.middleware.AccessForMomentAdmin)
		{
			typeFlnage.POST("/", h.createTypeFlange)
			typeFlnage.PUT("/:id", h.updateTypeFlange)
			typeFlnage.DELETE("/:id", h.deleteTypeFlange)
		}
	}
}

// @Summary Get Type Flange
// @Tags Sealur Moment -> type-flange
// @Security ApiKeyAuth
// @Description получение типов фланцев
// @ModuleID getTypeFlange
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{Data=[]moment_proto.TypeFlange}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/type-flange/ [get]
func (h *Handler) getTypeFlange(c *gin.Context) {
	tf, err := h.flangeClient.GetTypeFlange(c, &moment_proto.GetTypeFlangeRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: tf.TypeFlanges, Count: len(tf.TypeFlanges)})
}

// @Summary Create Type Flange
// @Tags Sealur Moment -> type-flange
// @Security ApiKeyAuth
// @Description создание тип фланца
// @ModuleID createTypeFlange
// @Accept json
// @Produce json
// @Param tf body models.TypeFlangeDTO true "type flange info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/type-flange/ [post]
func (h *Handler) createTypeFlange(c *gin.Context) {
	var dto models.TypeFlangeDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	tf, err := h.flangeClient.CreateTypeFlange(c, &moment_proto.CreateTypeFlangeRequest{Title: dto.Title})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Id: tf.Id, Message: "Created"})
}

// @Summary Update Type Flange
// @Tags Sealur Moment -> type-flange
// @Security ApiKeyAuth
// @Description обновление типа фланца
// @ModuleID updateTypeFlange
// @Accept json
// @Produce json
// @Param id path string true "type flange id"
// @Param tf body models.TypeFlangeDTO true "type flange info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/type-flange/{id} [put]
func (h *Handler) updateTypeFlange(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto models.TypeFlangeDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.flangeClient.UpdateTypeFlange(c, &moment_proto.UpdateTypeFlangeRequest{
		Id:    id,
		Title: dto.Title,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

// @Summary Delete Type Flange
// @Tags Sealur Moment -> type-flange
// @Security ApiKeyAuth
// @Description Удаление типа фланца
// @ModuleID deleteTypeFlange
// @Accept json
// @Produce json
// @Param id path string true "type flange id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/type-flange/{id} [delete]
func (h *Handler) deleteTypeFlange(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.flangeClient.DeleteTypeFlange(c, &moment_proto.DeleteTypeFlangeRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
