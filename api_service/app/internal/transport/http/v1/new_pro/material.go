package new_pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/snp_material_api"
	"github.com/gin-gonic/gin"
)

type MaterialHandler struct {
	materialApi snp_material_api.SnpMaterialServiceClient
}

func NewMaterialHandler(materialApi snp_material_api.SnpMaterialServiceClient) *MaterialHandler {
	return &MaterialHandler{
		materialApi: materialApi,
	}
}

func (h *Handler) initMaterialRoutes(api *gin.RouterGroup) *gin.RouterGroup {
	handler := NewMaterialHandler(h.materialApi)

	material := api.Group("/material")
	{
		material.GET("/", handler.get)
	}

	return material
}

func (h *MaterialHandler) get(c *gin.Context) {
	var dto *snp_material_api.GetSnpMaterial
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "отправлены некорректные данные")
		return
	}

	materials, err := h.materialApi.Get(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось получить материалы")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: materials})
}
