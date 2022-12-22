package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initFinningFactorRoutes(api *gin.RouterGroup) {
	finning := api.Group("/finning-factor", h.middleware.UserIdentity)
	{
		finning.GET("/", h.getFinningFactor)
		finning.POST("/", h.createFinningFactor)
		finning.POST("/few", h.createFewFinningFactor)
		finning.PUT("/:id", h.updateFinningFactor)
		finning.DELETE("/:id", h.deleteFinningFactor)
	}
}

func (h *Handler) getFinningFactor(c *gin.Context) {
	id := c.Query("devId")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	data, err := h.deviceClient.GetFinningFactor(c, &device_api.GetFinningFactorRequest{DevId: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: data.Finning, Count: len(data.Finning)})
}

func (h *Handler) createFinningFactor(c *gin.Context) {
	var dto moment_model.FinningFactorDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id, err := h.deviceClient.CreateFinningFactor(c, &device_api.CreateFinningFactorRequest{DevId: dto.DevId, Value: dto.Value})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Id: id.Id, Message: "Created"})
}

func (h *Handler) createFewFinningFactor(c *gin.Context) {
	var dto []moment_model.FinningFactorDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	finning := []*device_api.CreateFinningFactorRequest{}
	for _, d := range dto {
		finning = append(finning, &device_api.CreateFinningFactorRequest{DevId: d.DevId, Value: d.Value})
	}

	_, err := h.deviceClient.CreateFewFinningFactor(c, &device_api.CreateFewFinningFactorRequest{FinningFactor: finning})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

func (h *Handler) updateFinningFactor(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.FinningFactorDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.deviceClient.UpdateFinningFactor(c, &device_api.UpdateFinningFactorRequest{
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

func (h *Handler) deleteFinningFactor(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.deviceClient.DeleteFinningFactor(c, &device_api.DeleteFinningFactorRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
