package new_pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
	"github.com/Alexander272/sealur_proto/api/pro/ring_api"
	"github.com/Alexander272/sealur_proto/api/pro/rings_kit_api"
	"github.com/gin-gonic/gin"
)

type RingHandler struct {
	ringApi     ring_api.RingServiceClient
	ringsKitApi rings_kit_api.RingsKitServiceClient
	botApi      api.MostBotApi
}

func NewRingHandler(ringApi ring_api.RingServiceClient, ringsKitApi rings_kit_api.RingsKitServiceClient, botApi api.MostBotApi) *RingHandler {
	return &RingHandler{
		ringApi:     ringApi,
		ringsKitApi: ringsKitApi,
		botApi:      botApi,
	}
}

func (h *Handler) initRingsRoutes(api *gin.RouterGroup) {
	handler := NewRingHandler(h.ringApi, h.ringsKitApi, h.botApi)

	// rings := api.Group("/rings", h.middleware.UserIdentity)
	rings := api.Group("/rings")
	{
		rings.GET("single", handler.getSingle)
		rings.GET("kit", handler.getKit)
	}
}

func (h *RingHandler) getSingle(c *gin.Context) {
	ring, err := h.ringApi.Get(c, &ring_api.GetRings{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		h.botApi.SendError(c, err.Error(), "")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: ring})
}

func (h *RingHandler) getKit(c *gin.Context) {
	kit, err := h.ringsKitApi.Get(c, &rings_kit_api.GetRingsKit{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		h.botApi.SendError(c, err.Error(), "")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: kit})
}
