package new_pro

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/pro/snp_filler_api"
	"github.com/gin-gonic/gin"
)

type SnpFillerHandler struct {
	snpFillerApi snp_filler_api.SnpFillerServiceClient
	botApi       api.MostBotApi
}

func NewSnpFillerHandler(snpFillerApi snp_filler_api.SnpFillerServiceClient, botApi api.MostBotApi) *SnpFillerHandler {
	return &SnpFillerHandler{
		snpFillerApi: snpFillerApi,
		botApi:       botApi,
	}
}

func (h *Handler) initSnpFillerRoutes(api *gin.RouterGroup) {
	handler := NewSnpFillerHandler(h.snpFillerApi, h.botApi)

	// TODO проверять авторизацию
	filler := api.Group("/snp/fillers")
	{
		filler.GET("/", handler.get)
		// TODO только для админа
		filler.POST("/", handler.create)
		filler.POST("/several", handler.createSeveral)
		filler.PUT("/:id", handler.update)
		filler.DELETE("/:id", handler.delete)
	}
}

func (h *SnpFillerHandler) get(c *gin.Context) {
	standardId := c.Query("standardId")
	if standardId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "стандарт не задан")
		return
	}

	fillers, err := h.snpFillerApi.Get(c, &snp_filler_api.GetSnpFillers{StandardId: standardId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось получить наполнители")
		h.botApi.SendError(c, err.Error(), "")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: fillers})
}

func (h *SnpFillerHandler) create(c *gin.Context) {
	var dto *snp_filler_api.CreateSnpFiller
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	_, err := h.snpFillerApi.Create(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось создать наполнитель")

		body, bodyErr := json.Marshal(dto)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}
	c.JSON(http.StatusCreated, models.IdResponse{Message: "Наполнитель успешно создан"})
}

func (h *SnpFillerHandler) createSeveral(c *gin.Context) {
	var dto []*snp_filler_api.CreateSnpFiller
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	_, err := h.snpFillerApi.CreateSeveral(c, &snp_filler_api.CreateSeveralSnpFiller{SnpFillers: dto})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось создать наполнители")

		body, bodyErr := json.Marshal(dto)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}
	c.JSON(http.StatusCreated, models.IdResponse{Message: "Наполнители успешно созданы"})
}

func (h *SnpFillerHandler) update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "идентификатор не задан")
		return
	}

	var dto *snp_filler_api.UpdateSnpFiller
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}
	dto.Id = id

	_, err := h.snpFillerApi.Update(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось обновить наполнитель")

		body, bodyErr := json.Marshal(dto)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Наполнитель успешно обновлен"})
}

func (h *SnpFillerHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "идентификатор не задан")
		return
	}

	_, err := h.snpFillerApi.Delete(c, &snp_filler_api.DeleteSnpFiller{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось удалить наполнитель")
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "id": "%s" }`, id))
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Наполнитель успешно удален"})
}
