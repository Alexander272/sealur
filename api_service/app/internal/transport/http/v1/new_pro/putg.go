package new_pro

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/pro/putg_api"
	"github.com/gin-gonic/gin"
)

type PutgHandler struct {
	putgApi putg_api.PutgDataServiceClient
	botApi  api.MostBotApi
}

func NewPutgHandler(putgApi putg_api.PutgDataServiceClient, botApi api.MostBotApi) *PutgHandler {
	return &PutgHandler{
		putgApi: putgApi,
		botApi:  botApi,
	}
}

func (h *Handler) initPutgRoutes(api *gin.RouterGroup) {
	handler := NewPutgHandler(h.putgApi, h.botApi)

	// snp := api.Group("/snp", h.middleware.UserIdentity)
	putg := api.Group("/putg")
	{
		putg.GET("/base", handler.getBase)
		putg.GET("/data", handler.getData)
		putg.GET("", handler.get)
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

func (h *PutgHandler) getBase(c *gin.Context) {
	standardId := c.Query("standardId")
	emptyReq := c.Query("empty")
	empty, err := strconv.ParseBool(emptyReq)
	if err != nil {
		logger.Error("failed to parse empty. error: %w", err)
		empty = true
	}
	// typeFlangeId := c.Query("typeFlangeId")

	putg, err := h.putgApi.GetBase(c, &putg_api.GetPutgBase{StandardId: standardId, Empty: empty})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "standardId": "%s" }`, standardId))
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: putg})
}

func (h *PutgHandler) getData(c *gin.Context) {
	// стандарт на не стандартные фланцы
	standardId := c.Query("standardId")
	if standardId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param standardId", "стандарт не задан")
		return
	}
	constructionId := c.Query("constructionId")
	if constructionId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param constructionId", "конструкция не задана")
		return
	}
	baseConstructionId := c.Query("baseConstructionId")
	// if constructionId == "" {
	// 	models.NewErrorResponse(c, http.StatusBadRequest, "empty param base constructionId", "конструкция не задана")
	// 	return
	// }
	// fillerId := c.Query("fillerId")
	// configuration := c.Query("configuration")
	// changeStandard := c.Query("changeStandard")

	data, err := h.putgApi.GetData(c, &putg_api.GetPutgData{
		StandardId:         standardId,
		ConstructionId:     constructionId,
		BaseConstructionId: baseConstructionId,
		// Configuration:      configuration,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")

		body, bodyErr := json.Marshal(data)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: data})
}

func (h *PutgHandler) get(c *gin.Context) {
	fillerId := c.Query("fillerId")
	if fillerId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param fillerId", "материал прокладки не задан")
		return
	}
	baseId := c.Query("baseId")
	if baseId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param baseId", "материал прокладки не задан")
		return
	}
	flangeTypeId := c.Query("flangeTypeId")
	if flangeTypeId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param flangeTypeId", "тип фланца не задан")
		return
	}

	data, err := h.putgApi.Get(c, &putg_api.GetPutg{FillerId: fillerId, BaseId: baseId, FlangeTypeId: flangeTypeId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")

		body, bodyErr := json.Marshal(putg_api.GetPutg{FillerId: fillerId, BaseId: baseId, FlangeTypeId: flangeTypeId})
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: data})
}
