package new_pro

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/pro/snp_type_api"
	"github.com/gin-gonic/gin"
)

type SnpTypeHandler struct {
	snpTypeApi snp_type_api.SnpTypeServiceClient
	botApi     api.MostBotApi
}

func NewSnpTypeHandler(snpTypeApi snp_type_api.SnpTypeServiceClient, botApi api.MostBotApi) *SnpTypeHandler {
	return &SnpTypeHandler{
		snpTypeApi: snpTypeApi,
		botApi:     botApi,
	}
}

func (h *Handler) initSnpTypeRoutes(api *gin.RouterGroup) {
	handler := NewSnpTypeHandler(h.snpTypeApi, h.botApi)

	// TODO проверять авторизацию
	snpTypes := api.Group("snp/types")
	{
		snpTypes.GET("/", handler.get)
		// TODO только для админа
		snpTypes.POST("/", handler.create)
		snpTypes.PUT("/:id", handler.update)
		snpTypes.DELETE("/:id", handler.delete)
	}
}

func (h *SnpTypeHandler) get(c *gin.Context) {
	standardId := c.Query("standardId")
	if standardId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "стандарт не задан")
		return
	}

	snpTypes, err := h.snpTypeApi.Get(c, &snp_type_api.GetSnpTypes{StandardId: standardId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось получить типы")
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "standardId": "%s" }`, standardId))
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: snpTypes})
}

func (h *SnpTypeHandler) create(c *gin.Context) {
	var dto *snp_type_api.CreateSnpType
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	_, err := h.snpTypeApi.Create(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось создать тип")

		body, bodyErr := json.Marshal(dto)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}
	c.JSON(http.StatusCreated, models.IdResponse{Message: "Тип снп успешно создан"})
}

func (h *SnpTypeHandler) update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "идентификатор не задан")
		return
	}

	var dto *snp_type_api.UpdateSnpType
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}
	dto.Id = id

	_, err := h.snpTypeApi.Update(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось обновить тип")

		body, bodyErr := json.Marshal(dto)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Тип снп успешно обновлен"})
}

func (h *SnpTypeHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "идентификатор не задан")
		return
	}

	_, err := h.snpTypeApi.Delete(c, &snp_type_api.DeleteSnpType{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось удалить тип")
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "id": "%s" }`, id))
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Тип снп успешно удален"})
}
