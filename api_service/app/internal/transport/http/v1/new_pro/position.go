package new_pro

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/pro_model"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/file_api"
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
	"github.com/gin-gonic/gin"
)

type PositionHandler struct {
	positionApi position_api.PositionServiceClient
	fileApi     file_api.FileServiceClient
	botApi      api.MostBotApi
}

func NewPositionHandler(positionApi position_api.PositionServiceClient, fileApi file_api.FileServiceClient, botApi api.MostBotApi) *PositionHandler {
	return &PositionHandler{
		positionApi: positionApi,
		fileApi:     fileApi,
		botApi:      botApi,
	}
}

func (h *Handler) initPositionRoutes(api *gin.RouterGroup) {
	handler := NewPositionHandler(h.positionApi, h.fileApi, h.botApi)

	positions := api.Group("/positions", h.middleware.UserIdentity)
	{
		positions.POST("/", handler.create)
		positions.POST("/:id", handler.copy)
		positions.PUT("/:id", handler.update)
		positions.DELETE("/:id", handler.delete)
	}
}

func (h *PositionHandler) create(c *gin.Context) {
	var dto pro_model.Position
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные отправлены")
		return
	}

	position := dto.Parse()

	logger.Debug(position.PutgData)

	res, err := h.positionApi.Create(c, &position_api.CreatePosition{Position: position})
	if err != nil {
		if strings.Contains(err.Error(), "position exists") {
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Такая позиция уже добавлена в заявку")
			return
		}
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка")

		body, err := json.Marshal(dto)
		if err != nil {
			logger.Error("body error: ", err)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}

	// c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/orders/%s", order.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: res.Id, Message: "Created successfully"})
}

func (h *PositionHandler) copy(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty id param")
		return
	}

	var dto *position_api.CopyPosition
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные отправлены")
		return
	}
	dto.Id = id

	drawing, err := h.positionApi.Copy(c, dto)
	if err != nil {
		if strings.Contains(err.Error(), "position exists") {
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Такая позиция уже добавлена в заявку")
			return
		}
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка")

		body, err := json.Marshal(dto)
		if err != nil {
			logger.Error("body error: ", err)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}

	if drawing.Drawing != "" {
		parts := strings.Split(drawing.Drawing, "/")
		req := &file_api.CopyFileRequest{
			Id:       fmt.Sprintf("%s_%s", parts[len(parts)-2], parts[len(parts)-1]),
			Bucket:   "pro",
			Group:    dto.FromOrderId,
			NewGroup: dto.OrderId,
		}

		_, err := h.fileApi.Copy(c, req)
		if err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка при копировании чертежа")

			body, err := json.Marshal(req)
			if err != nil {
				logger.Error("body error: ", err)
			}
			h.botApi.SendError(c, err.Error(), string(body))

			return
		}
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Copied successfully"})
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

		body, err := json.Marshal(dto)
		if err != nil {
			logger.Error("body error: ", err)
		}
		h.botApi.SendError(c, err.Error(), string(body))

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
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "id": "%s" }`, id))
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Deleted successfully"})
}
