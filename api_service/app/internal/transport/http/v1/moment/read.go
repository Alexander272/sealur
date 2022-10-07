package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/read_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initReadRoutes(api *gin.RouterGroup) {
	read := api.Group("/data", h.middleware.UserIdentity)
	{
		read.GET("/flange", h.getFlange)
		read.GET("/float", h.getFloat)
	}
}

func (h *Handler) getFlange(c *gin.Context) {
	data, err := h.readClient.GetFlange(c, &read_api.GetFlangeRequest{})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: data})
}

func (h *Handler) getFloat(c *gin.Context) {
	data, err := h.readClient.GetFloat(c, &read_api.GetFloatRequest{})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: data})
}
