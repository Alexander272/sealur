package new_pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/pro_model"
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
	"github.com/gin-gonic/gin"
)

type PositionHandler struct {
	positionApi position_api.PositionServiceClient
}

func NewPositionHandler(positionApi position_api.PositionServiceClient) *PositionHandler {
	return &PositionHandler{
		positionApi: positionApi,
	}
}

func (h *Handler) initPositionRoutes(api *gin.RouterGroup) {
	handler := NewPositionHandler(h.positionApi)

	positions := api.Group("/positions", h.middleware.UserIdentity)
	{
		positions.POST("/", handler.create)
		positions.PUT("/:id", handler.update)
		positions.DELETE("/:id", handler.delete)
	}
}

func (h *PositionHandler) create(c *gin.Context) {
	var dto pro_model.Position
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	position := dto.Parse()

	res, err := h.positionApi.Create(c, &position_api.CreatePosition{Position: position})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	// c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/orders/%s", order.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: res.Id, Message: "Created successfully"})
}

func (h *PositionHandler) update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty id param")
		return
	}

	var dto pro_model.Position
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	position := dto.Parse()
	position.Id = id

	if _, err := h.positionApi.Update(c, &position_api.UpdatePosition{Position: position}); err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated successfully"})
}

func (h *PositionHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty id param")
		return
	}

	if _, err := h.positionApi.Delete(c, &position_api.DeletePosition{Id: id}); err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Deleted successfully"})
}
