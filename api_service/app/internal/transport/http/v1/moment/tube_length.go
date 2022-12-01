package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initTubeLengthRoutes(api *gin.RouterGroup) {
	tube := api.Group("/tube-length", h.middleware.UserIdentity)
	{
		tube.GET("/", h.getTubeLength)
		tube.POST("/", h.createTubeLength)
		tube.POST("/few", h.createFewTubeLength)
		tube.PUT("/:id", h.updateTubeLength)
		tube.DELETE("/:id", h.deleteTubeLength)
	}
}

func (h *Handler) getTubeLength(c *gin.Context) {
	id := c.Query("devId")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	data, err := h.deviceClient.GetTubeLength(c, &device_api.GetTubeLengthRequest{DevId: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: data.TubeLength, Count: len(data.TubeLength)})
}

func (h *Handler) createTubeLength(c *gin.Context) {
	var dto moment_model.TubeLengthDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id, err := h.deviceClient.CreateTubeLength(c, &device_api.CreateTubeLengthRequest{
		DevId: dto.DevId,
		Value: dto.Value,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Id: id.Id, Message: "Created"})
}

func (h *Handler) createFewTubeLength(c *gin.Context) {
	var dto []moment_model.TubeLengthDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	tube := []*device_api.CreateTubeLengthRequest{}
	for _, d := range dto {
		tube = append(tube, &device_api.CreateTubeLengthRequest{DevId: d.DevId, Value: d.Value})
	}

	_, err := h.deviceClient.CreateFewTubeLength(c, &device_api.CreateFewTubeLengthRequest{TubeLength: tube})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

func (h *Handler) updateTubeLength(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.TubeLengthDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.deviceClient.UpdateTubeLength(c, &device_api.UpdateTubeLengthRequest{
		Id:    id,
		DevId: dto.DevId,
		Value: dto.Value,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

func (h *Handler) deleteTubeLength(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.deviceClient.DeleteTubeLength(c, &device_api.DeleteTubeLengthRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Deleted"})
}
