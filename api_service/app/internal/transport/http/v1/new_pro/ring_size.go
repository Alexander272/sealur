package new_pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
	"github.com/Alexander272/sealur_proto/api/pro/ring_size_api"
	"github.com/gin-gonic/gin"
)

type RingSizeHandler struct {
	ringApi ring_size_api.RingSizeServiceClient
	botApi  api.MostBotApi
}

func NewRingSizeHandler(ringApi ring_size_api.RingSizeServiceClient, botApi api.MostBotApi) *RingSizeHandler {
	return &RingSizeHandler{
		ringApi: ringApi,
		botApi:  botApi,
	}
}

func (h *Handler) initRingSizeRoutes(api *gin.RouterGroup) {
	handler := NewRingSizeHandler(h.ringSizeApi, h.botApi)

	// rings := api.Group("/rings", h.middleware.UserIdentity)
	rings := api.Group("/ring-sizes")
	{
		rings.GET("", handler.get)
	}
}

func (h *RingSizeHandler) get(c *gin.Context) {
	ring, err := h.ringApi.GetAll(c, &ring_size_api.GetRingSize{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		h.botApi.SendError(c, err.Error(), "")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: ring})
}
