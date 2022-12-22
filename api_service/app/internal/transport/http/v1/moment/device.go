package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initDeviceRoutes(api *gin.RouterGroup) {
	device := api.Group("/device-mod", h.middleware.UserIdentity)
	{
		device.GET("/", h.getDevices)
		device.POST("/", h.createDevice)
		device.POST("/few", h.createDevices)
		device.PUT("/:id", h.updateDevice)
		device.DELETE("/:id", h.deleteDevice)
	}
}

func (h *Handler) getDevices(c *gin.Context) {
	device, err := h.deviceClient.GetDevice(c, &device_api.GetDeviceRequest{})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: device.Devices, Count: len(device.Devices)})
}

func (h *Handler) createDevice(c *gin.Context) {
	var dto moment_model.DeviceDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id, err := h.deviceClient.CreateDevice(c, &device_api.CreateDeviceRequest{Title: dto.Title})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Id: id.Id, Message: "Created"})
}

func (h *Handler) createDevices(c *gin.Context) {
	var dto []moment_model.DeviceDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	diveces := []*device_api.CreateDeviceRequest{}
	for _, d := range dto {
		diveces = append(diveces, &device_api.CreateDeviceRequest{Title: d.Title})
	}

	_, err := h.deviceClient.CreateFewDevice(c, &device_api.CreateFewDeviceRequest{Divices: diveces})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

func (h *Handler) updateDevice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.DeviceDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.deviceClient.UpdateDevice(c, &device_api.UpdateDeviceRequest{
		Id:    id,
		Title: dto.Title,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

func (h *Handler) deleteDevice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.deviceClient.DeleteDevice(c, &device_api.DeleteDeviceRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
