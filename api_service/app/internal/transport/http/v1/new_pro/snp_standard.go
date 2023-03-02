package new_pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/snp_standard_api"
	"github.com/gin-gonic/gin"
)

type SnpStandardHandler struct {
	snpStandardApi snp_standard_api.SnpStandardServiceClient
}

func NewSnpStandardHandler(snpStandardApi snp_standard_api.SnpStandardServiceClient) *SnpStandardHandler {
	return &SnpStandardHandler{
		snpStandardApi: snpStandardApi,
	}
}

func (h *Handler) initSnpStandardRoutes(api *gin.RouterGroup) {
	handler := NewSnpStandardHandler(h.snpStandardApi)

	// snp := api.Group("/snp", h.middleware.UserIdentity)
	standard := api.Group("/snp-standards")
	{
		standard.GET("/", handler.getAll)
	}
}

func (h *SnpStandardHandler) getAll(c *gin.Context) {
	standards, err := h.snpStandardApi.GetAll(c, &snp_standard_api.GetAllSnpStandards{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: standards.SnpStandards, Count: len(standards.SnpStandards)})
}
