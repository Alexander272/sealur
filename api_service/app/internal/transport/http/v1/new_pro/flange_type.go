package new_pro

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/pro/flange_type_api"
	"github.com/gin-gonic/gin"
)

type FlangeTypeHandler struct {
	flangeTypeApi flange_type_api.FlangeTypeServiceClient
	botApi        api.MostBotApi
}

func NewFlangeTypeHandler(flangeTypeApi flange_type_api.FlangeTypeServiceClient, botApi api.MostBotApi) *FlangeTypeHandler {
	return &FlangeTypeHandler{
		flangeTypeApi: flangeTypeApi,
		botApi:        botApi,
	}
}

func (h *Handler) initFlangeTypeRoutes(api *gin.RouterGroup) {
	handler := NewFlangeTypeHandler(h.flangeTypeApi, h.botApi)

	// TODO проверять авторизацию
	flangeType := api.Group("flange-types")
	{
		flangeType.GET("/", handler.get)
		// TODO только для админа
		flangeType.POST("/", handler.create)
		flangeType.PUT("/:id", handler.update)
		flangeType.DELETE("/:id", handler.delete)
	}
}

func (h *FlangeTypeHandler) get(c *gin.Context) {
	flangeTypes, err := h.flangeTypeApi.Get(c, &flange_type_api.GetFlangeType{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось получить типы фланцев")
		h.botApi.SendError(c, err.Error(), "")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: flangeTypes})
}

func (h *FlangeTypeHandler) create(c *gin.Context) {
	var dto *flange_type_api.CreateFlangeType
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	_, err := h.flangeTypeApi.Create(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось создать тип фланца")

		body, err := json.Marshal(dto)
		if err != nil {
			logger.Error("body error: ", err)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}
	c.JSON(http.StatusCreated, models.IdResponse{Message: "Тип фланца успешно создан"})
}

func (h *FlangeTypeHandler) update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty id param")
		return
	}

	var dto *flange_type_api.UpdateFlangeType
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "отправлены некорректные данные")
		return
	}
	dto.Id = id

	_, err := h.flangeTypeApi.Update(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось обновить тип фланца")

		body, err := json.Marshal(dto)
		if err != nil {
			logger.Error("body error: ", err)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Тип фланца успешно обновлен"})
}

func (h *FlangeTypeHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty id param")
		return
	}

	_, err := h.flangeTypeApi.Delete(c, &flange_type_api.DeleteFlangeType{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось удалить тип фланца")
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "id": "%s" }`, id))
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Тип фланца успешно удален"})
}
