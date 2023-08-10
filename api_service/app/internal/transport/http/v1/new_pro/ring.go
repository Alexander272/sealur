package new_pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
	"github.com/Alexander272/sealur_proto/api/pro/ring_api"
	"github.com/gin-gonic/gin"
)

type RingHandler struct {
	ringApi ring_api.RingServiceClient
	botApi  api.MostBotApi
}

func NewRingHandler(ringApi ring_api.RingServiceClient, botApi api.MostBotApi) *RingHandler {
	return &RingHandler{
		ringApi: ringApi,
		botApi:  botApi,
	}
}

func (h *Handler) initRingsRoutes(api *gin.RouterGroup) {
	handler := NewRingHandler(h.ringApi, h.botApi)

	// rings := api.Group("/rings", h.middleware.UserIdentity)
	rings := api.Group("/rings")
	{
		rings.GET("", handler.get)
	}
}

func (h *RingHandler) get(c *gin.Context) {
	ring, err := h.ringApi.Get(c, &ring_api.GetRings{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		h.botApi.SendError(c, err.Error(), "")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: ring})
}
