package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initNumberOfMovesRoutes(api *gin.RouterGroup) {
	number := api.Group("/number-of-moves", h.middleware.UserIdentity)
	{
		number.GET("/", h.getNumberOfMoves)
		number.POST("/", h.createNumberOfMoves)
		number.POST("/few", h.createFewNumberOfMoves)
		number.PUT("/:id", h.updateNumberOfMoves)
		number.DELETE("/:id", h.deleteNumberOfMoves)
	}
}

func (h *Handler) getNumberOfMoves(c *gin.Context) {
	id := c.Query("devId")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	data, err := h.deviceClient.GetNumberOfMoves(c, &device_api.GetNumberOfMovesRequest{DevId: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: data.Number, Count: len(data.Number)})
}

func (h *Handler) createNumberOfMoves(c *gin.Context) {
	var dto moment_model.NumberOfMovesDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id, err := h.deviceClient.CreateNumberOfMoves(c, &device_api.CreateNumberOfMovesRequest{
		DevId:   dto.DevId,
		CountId: dto.CountId,
		Value:   dto.Value,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Id: id.Id, Message: "Created"})
}

func (h *Handler) createFewNumberOfMoves(c *gin.Context) {
	var dto []moment_model.NumberOfMovesDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	numbers := []*device_api.CreateNumberOfMovesRequest{}
	for _, d := range dto {
		numbers = append(numbers, &device_api.CreateNumberOfMovesRequest{
			DevId:   d.DevId,
			CountId: d.CountId,
			Value:   d.Value,
		})
	}

	_, err := h.deviceClient.CreateFewNumberOfMoves(c, &device_api.CreateFewNumberOfMovesRequest{Number: numbers})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

func (h *Handler) updateNumberOfMoves(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.NumberOfMovesDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.deviceClient.UpdateNumberOfMoves(c, &device_api.UpdateNumberOfMovesRequest{
		Id:      id,
		DevId:   dto.DevId,
		CountId: dto.CountId,
		Value:   dto.Value,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

func (h *Handler) deleteNumberOfMoves(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.deviceClient.DeleteNumberOfMoves(c, &device_api.DeleteNumberOfMovesRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Deleted"})
}
