package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initReadRoutes(api *gin.RouterGroup) {
	read := api.Group("/default", h.middleware.UserIdentity)
	{
		read.GET("/flange", h.getDefFlange)
	}
}

func (h *Handler) getDefFlange(c *gin.Context) {
	data, err := h.readClient.GetFlange(c, &moment_api.GetFlangeRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: data})
}
