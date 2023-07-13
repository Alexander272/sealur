package new_pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_size_api"
	"github.com/gin-gonic/gin"
)

type PutgSizeHandler struct {
	sizeApi putg_size_api.PutgSizeServiceClient
	botApi  api.MostBotApi
}

func NewPutgSizeHandler(sizeApi putg_size_api.PutgSizeServiceClient, botApi api.MostBotApi) *PutgSizeHandler {
	return &PutgSizeHandler{
		sizeApi: sizeApi,
		botApi:  botApi,
	}
}

func (h *Handler) initPutgSizeRoutes(api *gin.RouterGroup) {
	handler := NewPutgSizeHandler(h.putgSizeApi, h.botApi)

	// TODO проверять авторизацию
	sizes := api.Group("/putg/sizes")
	{
		sizes.GET("", handler.get)
		// TODO только для админа
		// sizes.POST("/", handler.create)
		// sizes.POST("/several", handler.createSeveral)
		// sizes.PUT("/:id", handler.update)
		// sizes.DELETE("/:id", handler.delete)
	}
}

func (h *PutgSizeHandler) get(c *gin.Context) {
	flangeTypeId := c.Query("flangeTypeId")
	if flangeTypeId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "тип фланца не задан")
		return
	}
	baseConstructionId := c.Query("baseConstructionId")
	if baseConstructionId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "тип конструкции не задан")
		return
	}
	baseFillerId := c.Query("baseFillerId")
	if baseFillerId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "тип наполнителя не задан")
		return
	}

	sizes, err := h.sizeApi.Get(c, &putg_size_api.GetPutgSize{FlangeTypeId: flangeTypeId, BaseConstructionId: baseConstructionId, BaseFillerId: baseFillerId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось получить размеры")
		h.botApi.SendError(c, err.Error(), "")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: sizes})
}
