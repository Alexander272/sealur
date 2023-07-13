package new_pro

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/pro/snp_api"
	"github.com/gin-gonic/gin"
)

type SnpHandler struct {
	snpApi snp_api.SnpDataServiceClient
	botApi api.MostBotApi
}

func NewSnpHandler(snpApi snp_api.SnpDataServiceClient, botApi api.MostBotApi) *SnpHandler {
	return &SnpHandler{
		snpApi: snpApi,
		botApi: botApi,
	}
}

func (h *Handler) initSNPRoutes(api *gin.RouterGroup) {
	handler := NewSnpHandler(h.snpApi, h.botApi)

	// snp := api.Group("/snp", h.middleware.UserIdentity)
	snp := api.Group("/snp-new")
	{
		snp.GET("", handler.get)
		snp.GET("/data", handler.getData)
		// snp.GET("/default", h.getDefault)
		// snp.GET("/", h.getSNP)
		// snp = snp.Group("/", h.middleware.AccessForProAdmin)
		// {
		// 	snp.POST("/", h.createSNP)
		// 	snp.PUT("/:id", h.updateSNP)
		// 	snp.DELETE("/:id", h.deleteSNP)
		// }
	}
}

func (h *SnpHandler) get(c *gin.Context) {
	// standardId := c.Query("standardId")
	// if standardId == "" {
	// 	models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty standard id param")
	// 	return
	// }
	// snpStandardId := c.Query("snpStandardId")
	// if snpStandardId == "" {
	// 	models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty snp standard id param")
	// 	return
	// }
	snpTypeId := c.Query("typeId")
	if snpTypeId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty snp type id param")
		return
	}
	// flangeId := c.Query("flangeId")
	// if flangeId == "" {
	// 	models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty flange standard id param")
	// 	return
	// }
	hasD2Query := c.Query("hasD2")
	hasD2 := false
	if hasD2Query != "" {
		var err error
		hasD2, err = strconv.ParseBool(hasD2Query)
		if err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty has d2 param")
			return
		}
	}

	snp, err := h.snpApi.Get(c, &snp_api.GetSnp{
		// StandardId:       standardId,
		// SnpStandardId:    snpStandardId,
		SnpTypeId: snpTypeId,
		// FlangeStandardId: flangeId,
		HasD2: hasD2,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")

		body, bodyErr := json.Marshal(snp)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: snp})
}

func (h *SnpHandler) getData(c *gin.Context) {
	standardId := c.Query("standardId")
	// if standardId == "" {
	// 	models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty standard id param")
	// 	return
	// }
	snpStandardId := c.Query("snpStandardId")

	snp, err := h.snpApi.GetData(c, &snp_api.GetSnpData{
		StandardId:    standardId,
		SnpStandardId: snpStandardId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")

		body, bodyErr := json.Marshal(snp)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: snp.SnpData})
}
