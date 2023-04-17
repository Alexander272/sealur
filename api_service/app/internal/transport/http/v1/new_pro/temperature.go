package new_pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/temperature_api"
	"github.com/gin-gonic/gin"
)

type TemperatureHandler struct {
	temperatureApi temperature_api.TemperatureServiceClient
}

func NewTemperatureHandler(temperatureApi temperature_api.TemperatureServiceClient) *TemperatureHandler {
	return &TemperatureHandler{
		temperatureApi: temperatureApi,
	}
}

func (h *Handler) initTemperatureRoutes(api *gin.RouterGroup) {
	handler := NewTemperatureHandler(h.temperatureApi)

	// TODO проверять авторизацию
	temperature := api.Group("/temperature")
	{
		temperature.GET("/", handler.get)
		// TODO только для админа
		temperature.POST("/", handler.create)
		temperature.PUT("/:id", handler.update)
		temperature.DELETE("/:id", handler.delete)
	}
}

func (h *TemperatureHandler) get(c *gin.Context) {
	temperature, err := h.temperatureApi.GetAll(c, &temperature_api.GetAllTemperatures{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось получить температуры")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: temperature})
}

func (h *TemperatureHandler) create(c *gin.Context) {
	var dto *temperature_api.CreateTemperature
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	_, err := h.temperatureApi.Create(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось создать температуру")
		return
	}
	c.JSON(http.StatusCreated, models.IdResponse{Message: "Температура успешно создана"})
}

func (h *TemperatureHandler) update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "идентификатор не задан")
		return
	}

	var dto *temperature_api.UpdateTemperature
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}
	dto.Id = id

	_, err := h.temperatureApi.Update(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось обновить температуру")
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Температура успешно обновлена"})
}

func (h *TemperatureHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "идентификатор не задан")
		return
	}

	_, err := h.temperatureApi.Delete(c, &temperature_api.DeleteTemperature{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось удалить температуру")
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Температура успешно удалена"})
}
