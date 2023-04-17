package new_pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/snp_data_api"
	"github.com/gin-gonic/gin"
)

type SnpDataHandler struct {
	dataApi snp_data_api.SnpDataServiceClient
}

func NewSnpDataHandler(dataApi snp_data_api.SnpDataServiceClient) *SnpDataHandler {
	return &SnpDataHandler{
		dataApi: dataApi,
	}
}

func (h *Handler) initSnpDataRoutes(api *gin.RouterGroup) {
	handler := NewSnpDataHandler(h.snpDataApi)

	// TODO проверять авторизацию
	data := api.Group("/snp/data")
	{
		data.GET("/", handler.get)
		// TODO только для админа
		data.POST("/", handler.create)
		data.PUT("/:id", handler.update)
		data.DELETE("/:id", handler.delete)
	}
}

func (h *SnpDataHandler) get(c *gin.Context) {
	typeId := c.Query("typeId")
	if typeId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "тип снп не задан")
		return
	}

	data, err := h.dataApi.Get(c, &snp_data_api.GetSnpData{TypeId: typeId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось получить данные о прокладке")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: data})
}

func (h *SnpDataHandler) create(c *gin.Context) {
	var dto *snp_data_api.CreateSnpData
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	_, err := h.dataApi.Create(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось добавить данные о прокладке")
		return
	}
	c.JSON(http.StatusCreated, models.IdResponse{Message: "Данные о прокладке успешно добавлены"})
}

func (h *SnpDataHandler) update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "идентификатор не задан")
		return
	}

	var dto *snp_data_api.UpdateSnpData
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}
	dto.Id = id

	_, err := h.dataApi.Update(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось обновить данные о прокладке")
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Данные о прокладке успешно обновлены"})
}

func (h *SnpDataHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "идентификатор не задан")
		return
	}

	_, err := h.dataApi.Delete(c, &snp_data_api.DeleteSnpData{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось удалить данные о прокладке")
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Данные о прокладке успешно удален"})
}
