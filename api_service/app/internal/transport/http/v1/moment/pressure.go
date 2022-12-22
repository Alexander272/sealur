package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initPressureRoutes(api *gin.RouterGroup) {
	pressure := api.Group("/pressure", h.middleware.UserIdentity)
	{
		pressure.GET("/", h.getPressure)
		pressure.POST("/", h.createPressure)
		pressure.POST("/few", h.createFewPressure)
		pressure.PUT("/:id", h.updatePressure)
		pressure.DELETE("/:id", h.deletePressure)
	}
}

func (h *Handler) getPressure(c *gin.Context) {
	pres, err := h.deviceClient.GetPressure(c, &device_api.GetPressureRequest{})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: pres.Pressures, Count: len(pres.Pressures)})
}

func (h *Handler) createPressure(c *gin.Context) {
	var dto moment_model.PressureDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id, err := h.deviceClient.CreatePressure(c, &device_api.CreatePressureRequest{Value: dto.Value})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Id: id.Id, Message: "Created"})
}

func (h *Handler) createFewPressure(c *gin.Context) {
	var dto []moment_model.PressureDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	pres := []*device_api.CreatePressureRequest{}
	for _, d := range dto {
		pres = append(pres, &device_api.CreatePressureRequest{Value: d.Value})
	}

	_, err := h.deviceClient.CreateFewPressure(c, &device_api.CreateFewPressureRequest{Pressure: pres})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

func (h *Handler) updatePressure(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.PressureDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.deviceClient.UpdatePressure(c, &device_api.UpdatePressureRequest{
		Id:    id,
		Value: dto.Value,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

func (h *Handler) deletePressure(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.deviceClient.DeletePressure(c, &device_api.DeletePressureRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
