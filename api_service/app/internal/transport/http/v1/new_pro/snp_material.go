package new_pro

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/pro/snp_material_api"
	"github.com/gin-gonic/gin"
)

type SnpMaterialHandler struct {
	materialApi snp_material_api.SnpMaterialServiceClient
	botApi      api.MostBotApi
}

func NewSnpMaterialHandler(materialApi snp_material_api.SnpMaterialServiceClient, botApi api.MostBotApi) *SnpMaterialHandler {
	return &SnpMaterialHandler{
		materialApi: materialApi,
		botApi:      botApi,
	}
}

func (h *Handler) initSnpMaterialRoutes(api *gin.RouterGroup) *gin.RouterGroup {
	handler := NewSnpMaterialHandler(h.snpMaterialApi, h.botApi)

	// TODO проверять авторизацию
	material := api.Group("/snp/materials")
	{
		material.GET("/", handler.get)
		// TODO только для админа
		material.POST("/", handler.create)
		material.PUT("/:id", handler.update)
		material.DELETE("/:id", handler.delete)
	}

	return material
}

func (h *SnpMaterialHandler) get(c *gin.Context) {
	standardId := c.Query("standardId")
	if standardId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "стандарт не задан")
		return
	}

	materials, err := h.materialApi.Get(c, &snp_material_api.GetSnpMaterial{StandardId: standardId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось получить материалы для стандарта")
		h.botApi.SendError(c, err.Error(), "")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: materials})
}

func (h *SnpMaterialHandler) create(c *gin.Context) {
	var dto *snp_material_api.CreateSnpMaterial
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "отправлены некорректные данные")
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

func (h *SnpMaterialHandler) update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty id param")
		return
	}

	var dto *snp_material_api.UpdateSnpMaterial
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "отправлены некорректные данные")

		body, bodyErr := json.Marshal(dto)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}
	dto.Id = id

	_, err := h.materialApi.Update(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось обновить материал")
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Материал успешно обновлен"})
}

func (h *SnpMaterialHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty id param")
		return
	}

	_, err := h.materialApi.Delete(c, &snp_material_api.DeleteSnpMaterial{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось удалить материал")
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "id": "%s" }`, id))
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Материал успешно удален"})
}
