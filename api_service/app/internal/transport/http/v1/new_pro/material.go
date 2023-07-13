package new_pro

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/pro/material_api"
	"github.com/gin-gonic/gin"
)

type MaterialHandler struct {
	materialApi material_api.MaterialServiceClient
	botApi      api.MostBotApi
}

func NewMaterialHandler(materialApi material_api.MaterialServiceClient, botApi api.MostBotApi) *MaterialHandler {
	return &MaterialHandler{
		materialApi: materialApi,
		botApi:      botApi,
	}
}

func (h *Handler) initMaterialRoutes(api *gin.RouterGroup) {
	handler := NewMaterialHandler(h.materialApi, h.botApi)

	// TODO проверять авторизацию
	materials := api.Group("/materials")
	{
		materials.GET("/", handler.get)
		// TODO только для админа
		materials.POST("/", handler.create)
		materials.PUT("/:id", handler.update)
		materials.DELETE("/:id", handler.delete)
	}
}

func (h *MaterialHandler) get(c *gin.Context) {
	materials, err := h.materialApi.GetAll(c, &material_api.GetAllMaterials{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось получить материалы")
		h.botApi.SendError(c, err.Error(), "")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: materials})
}

func (h *MaterialHandler) create(c *gin.Context) {
	var dto *material_api.CreateMaterial
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	_, err := h.materialApi.Create(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось создать материал")

		body, bodyErr := json.Marshal(dto)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}
	c.JSON(http.StatusCreated, models.IdResponse{Message: "Материал успешно создан"})
}

func (h *MaterialHandler) update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty id param")
		return
	}

	var dto *material_api.UpdateMaterial
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "отправлены некорректные данные")
		return
	}
	dto.Id = id

	_, err := h.materialApi.Update(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось обновить материал")

		body, bodyErr := json.Marshal(dto)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Материал успешно обновлен"})
}

func (h *MaterialHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty id param")
		return
	}

	_, err := h.materialApi.Delete(c, &material_api.DeleteMaterial{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось удалить материал")
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "id": "%s" }`, id))
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Материал успешно удален"})
}
