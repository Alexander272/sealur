package new_pro

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/pro/snp_size_api"
	"github.com/gin-gonic/gin"
)

type SizeHandler struct {
	sizeApi snp_size_api.SnpSizeServiceClient
	botApi  api.MostBotApi
}

func NewSizeHandler(sizeApi snp_size_api.SnpSizeServiceClient, botApi api.MostBotApi) *SizeHandler {
	return &SizeHandler{
		sizeApi: sizeApi,
		botApi:  botApi,
	}
}

func (h *Handler) initSizeRoutes(api *gin.RouterGroup) {
	handler := NewSizeHandler(h.sizeApi, h.botApi)

	// TODO проверять авторизацию
	sizes := api.Group("/snp/sizes")
	{
		sizes.GET("/", handler.get)
		// TODO только для админа
		sizes.POST("/", handler.create)
		sizes.POST("/several", handler.createSeveral)
		sizes.PUT("/:id", handler.update)
		sizes.DELETE("/:id", handler.delete)
	}
}

func (h *SizeHandler) get(c *gin.Context) {
	typeId := c.Query("typeId")
	if typeId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "тип снп не задан")
		return
	}

	sizes, err := h.sizeApi.Get(c, &snp_size_api.GetSnpSize{TypeId: typeId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось получить размеры")
		h.botApi.SendError(c, err.Error(), "")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: sizes})
}

func (h *SizeHandler) create(c *gin.Context) {
	var dto *snp_size_api.CreateSnpSize
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	_, err := h.sizeApi.Create(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось создать размеры")

		body, err := json.Marshal(dto)
		if err != nil {
			logger.Error("body error: ", err)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}
	c.JSON(http.StatusCreated, models.IdResponse{Message: "Размеры успешно созданы"})
}

func (h *SizeHandler) createSeveral(c *gin.Context) {
	var dto []*snp_size_api.CreateSnpSize
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	_, err := h.sizeApi.CreateSeveral(c, &snp_size_api.CreateSeveralSnpSize{Sizes: dto})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось создать размеры")

		body, err := json.Marshal(dto)
		if err != nil {
			logger.Error("body error: ", err)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}
	c.JSON(http.StatusCreated, models.IdResponse{Message: "Размеры успешно созданы"})
}

func (h *SizeHandler) update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "идентификатор не задан")
		return
	}

	var dto *snp_size_api.UpdateSnpSize
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}
	dto.Id = id

	_, err := h.sizeApi.Update(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось обновить размеры")

		body, err := json.Marshal(dto)
		if err != nil {
			logger.Error("body error: ", err)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "размеры успешно обновлены"})
}

func (h *SizeHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "идентификатор не задан")
		return
	}

	_, err := h.sizeApi.Delete(c, &snp_size_api.DeleteSnpSize{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось удалить размеры")
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "id": "%s" }`, id))
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "размеры успешно удалены"})
}
