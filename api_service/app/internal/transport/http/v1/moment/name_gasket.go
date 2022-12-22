package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initNameGasketRoutes(api *gin.RouterGroup) {
	gasket := api.Group("/name-gasket", h.middleware.UserIdentity)
	{
		gasket.GET("/", h.getNameGasket)
		gasket.GET("/full", h.getFullNameGasket)
		gasket.GET("/size/:id", h.getNameGasketSize)
		gasket.POST("/", h.createNameGasket)
		gasket.POST("/few", h.createFewNameGasket)
		gasket.PUT("/:id", h.updateNameGasket)
		gasket.DELETE("/:id", h.deleteNameGasket)
	}
}

func (h *Handler) getNameGasket(c *gin.Context) {
	finId := c.Query("finId")
	if finId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty finId", "empty finIf param")
		return
	}
	presId := c.Query("presId")
	numId := c.Query("numId")

	data, err := h.deviceClient.GetNameGasket(c, &device_api.GetNameGasketRequest{FinId: finId, PresId: presId, NumId: numId})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: data.Gasket, Count: len(data.Gasket)})
}

func (h *Handler) getFullNameGasket(c *gin.Context) {
	id := c.Query("finId")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty finId", "empty finId param")
		return
	}
	presId := c.Query("presId")
	numId := c.Query("numId")

	data, err := h.deviceClient.GetFullNameGasket(c, &device_api.GetFullNameGasketRequest{FinId: id, PresId: presId, NumId: numId})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: data.Gasket, Count: len(data.Gasket)})
}

func (h *Handler) getNameGasketSize(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	data, err := h.deviceClient.GetNameGasketSize(c, &device_api.GetNameGasketSizeRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: data.Gasket, Count: len(data.Gasket)})
}

func (h *Handler) createNameGasket(c *gin.Context) {
	var dto moment_model.NameGasketDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	gasket := &device_api.CreateNameGasketRequest{
		FinId:     dto.FinId,
		NumId:     dto.NumId,
		PresId:    dto.PresId,
		Title:     dto.Title,
		SizeLong:  dto.SizeLong,
		SizeTrans: dto.SizeTrans,
		Width:     dto.Width,
		Thick1:    dto.Thick1,
		Thick2:    dto.Thick2,
		Thick3:    dto.Thick3,
		Thick4:    dto.Thick4,
	}

	id, err := h.deviceClient.CreateNameGasket(c, gasket)
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Id: id.Id, Message: "Created"})
}

func (h *Handler) createFewNameGasket(c *gin.Context) {
	var dto []moment_model.NameGasketDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	gaskets := []*device_api.CreateNameGasketRequest{}
	for _, d := range dto {
		gaskets = append(gaskets, &device_api.CreateNameGasketRequest{
			FinId:     d.FinId,
			NumId:     d.NumId,
			PresId:    d.PresId,
			Title:     d.Title,
			SizeLong:  d.SizeLong,
			SizeTrans: d.SizeTrans,
			Width:     d.Width,
			Thick1:    d.Thick1,
			Thick2:    d.Thick2,
			Thick3:    d.Thick3,
			Thick4:    d.Thick4,
		})
	}

	_, err := h.deviceClient.CreateFewNameGasket(c, &device_api.CreateFewNameGasketRequest{Gasket: gaskets})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusBadRequest, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

func (h *Handler) updateNameGasket(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.NameGasketDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	gasket := &device_api.UpdateNameGasketRequest{
		Id:        id,
		FinId:     dto.FinId,
		NumId:     dto.NumId,
		PresId:    dto.PresId,
		Title:     dto.Title,
		SizeLong:  dto.SizeLong,
		SizeTrans: dto.SizeTrans,
		Width:     dto.Width,
		Thick1:    dto.Thick1,
		Thick2:    dto.Thick2,
		Thick3:    dto.Thick3,
		Thick4:    dto.Thick4,
	}

	_, err := h.deviceClient.UpdateNameGasket(c, gasket)
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

func (h *Handler) deleteNameGasket(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.deviceClient.DeleteNameGasket(c, &device_api.DeleteNameGasketRequest{Id: id})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Deleted"})
}
