package new_pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
	"github.com/Alexander272/sealur_proto/api/pro/ring_modifying_api"
	"github.com/gin-gonic/gin"
)

type RingModifyingHandler struct {
	ringApi ring_modifying_api.RingModifyingServiceClient
	botApi  api.MostBotApi
}

func NewRingModifyingHandler(ringApi ring_modifying_api.RingModifyingServiceClient, botApi api.MostBotApi) *RingModifyingHandler {
	return &RingModifyingHandler{
		ringApi: ringApi,
		botApi:  botApi,
	}
}

func (h *Handler) initRingModifyingRoutes(api *gin.RouterGroup) {
	handler := NewRingModifyingHandler(h.ringModifyingApi, h.botApi)

	// rings := api.Group("/rings", h.middleware.UserIdentity)
	rings := api.Group("/ring-modifying")
	{
		rings.GET("", handler.get)
	}
}

func (h *RingModifyingHandler) get(c *gin.Context) {
	ring, err := h.ringApi.GetAll(c, &ring_modifying_api.GetRingModifying{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		h.botApi.SendError(c, err.Error(), "")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: ring})
}
