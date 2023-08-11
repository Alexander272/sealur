package new_pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
	"github.com/Alexander272/sealur_proto/api/pro/ring_material_api"
	"github.com/gin-gonic/gin"
)

type RingMaterialHandler struct {
	ringApi ring_material_api.RingMaterialServiceClient
	botApi  api.MostBotApi
}

func NewRingMaterialHandler(ringApi ring_material_api.RingMaterialServiceClient, botApi api.MostBotApi) *RingMaterialHandler {
	return &RingMaterialHandler{
		ringApi: ringApi,
		botApi:  botApi,
	}
}

func (h *Handler) initRingMaterialsRoutes(api *gin.RouterGroup) {
	handler := NewRingMaterialHandler(h.ringMaterialApi, h.botApi)

	// rings := api.Group("/rings", h.middleware.UserIdentity)
	rings := api.Group("/ring-materials")
	{
		rings.GET("/:type", handler.get)
	}
}

func (h *RingMaterialHandler) get(c *gin.Context) {
	typeMaterial := c.Param("type")
	if typeMaterial == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty type", "тип материалов не задан")
		return
	}

	ring, err := h.ringApi.Get(c, &ring_material_api.GetRingMaterial{Type: typeMaterial})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		h.botApi.SendError(c, err.Error(), "")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: ring})
}
