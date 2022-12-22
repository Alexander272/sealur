package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initSectionExecutionRoutes(api *gin.RouterGroup) {
	section := api.Group("section-execution", h.middleware.UserIdentity)
	{
		section.GET("/", h.getSectionExecution)
		section.POST("/", h.createSectionExecution)
		section.POST("/few", h.createFewSectionExecution)
		section.PUT("/:id", h.updateSectionExecution)
		section.DELETE("/:id", h.deleteSectionExecution)
	}
}

func (h *Handler) getSectionExecution(c *gin.Context) {
	id := c.Query("devId")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	data, err := h.deviceClient.GetSectionExecution(c, &device_api.GetSectionExecutionRequest{DevId: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: data.Section, Count: len(data.Section)})
}

func (h *Handler) createSectionExecution(c *gin.Context) {
	var dto moment_model.SectionExecutionDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id, err := h.deviceClient.CreateSectionExecution(c, &device_api.CreateSectionExecutionRequest{DevId: dto.DevId, Value: dto.Value})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Id: id.Id, Message: "Created"})
}

func (h *Handler) createFewSectionExecution(c *gin.Context) {
	var dto []moment_model.SectionExecutionDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	section := []*device_api.CreateSectionExecutionRequest{}
	for _, d := range dto {
		section = append(section, &device_api.CreateSectionExecutionRequest{DevId: d.DevId, Value: d.Value})
	}

	_, err := h.deviceClient.CreateFewSectionExecution(c, &device_api.CreateFewSectionExecutionRequest{Section: section})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

func (h *Handler) updateSectionExecution(c *gin.Context) {
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

	_, err := h.deviceClient.UpdateSectionExecution(c, &device_api.UpdateSectionExecutionRequest{
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

func (h *Handler) deleteSectionExecution(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.deviceClient.DeleteSectionExecution(c, &device_api.DeleteSectionExecutionRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
