package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initTubeCountRoutes(api *gin.RouterGroup) {
	tube := api.Group("/tube-count", h.middleware.UserIdentity)
	{
		tube.GET("/", h.getTubeCount)
		tube.POST("/", h.createTubeCount)
		tube.POST("/few", h.createFewTubeCount)
		tube.PUT("/:id", h.updateTubeCount)
		tube.DELETE("/:id", h.deleteTubeCount)
	}
}

func (h *Handler) getTubeCount(c *gin.Context) {
	tube, err := h.deviceClient.GetTubeCount(c, &device_api.GetTubeCountRequest{})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: tube.TubeCount, Count: len(tube.TubeCount)})
}

func (h *Handler) createTubeCount(c *gin.Context) {
	var dto moment_model.TubeCountDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id, err := h.deviceClient.CreateTubeCount(c, &device_api.CreateTubeCountRequest{Value: dto.Value})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Id: id.Id, Message: "Created"})
}

func (h *Handler) createFewTubeCount(c *gin.Context) {
	var dto []moment_model.TubeCountDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	tube := []*device_api.CreateTubeCountRequest{}
	for _, d := range dto {
		tube = append(tube, &device_api.CreateTubeCountRequest{Value: d.Value})
	}

	_, err := h.deviceClient.CreateFewTubeCount(c, &device_api.CreateFewTubeCountRequest{TubeCount: tube})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

func (h *Handler) updateTubeCount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.TubeCountDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.deviceClient.UpdateTubeCount(c, &device_api.UpdateTubeCountRequest{
		Id:    id,
		Value: dto.Value,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

func (h *Handler) deleteTubeCount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.deviceClient.DeleteTubeCount(c, &device_api.DeleteTubeCountRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
