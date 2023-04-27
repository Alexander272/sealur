package new_pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/putg_api"
	"github.com/gin-gonic/gin"
)

type PutgHandler struct {
	putgApi putg_api.PutgDataServiceClient
}

func NewPutgHandler(putgApi putg_api.PutgDataServiceClient) *PutgHandler {
	return &PutgHandler{
		putgApi: putgApi,
	}
}

func (h *Handler) initPutgRoutes(api *gin.RouterGroup) {
	handler := NewPutgHandler(h.putgApi)

	// snp := api.Group("/snp", h.middleware.UserIdentity)
	putg := api.Group("/putg")
	{
		putg.GET("/base", handler.getBase)
		putg.GET("/data", handler.getData)
		// putg.GET("", handler.get)
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
	putg, err := h.putgApi.GetBase(c, &putg_api.GetPutgBase{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: putg})
}

func (h *PutgHandler) getData(c *gin.Context) {
	//TODO стандарт на не стандартные фланцы
	standardId := c.Query("standardId")
	// if standardId == "" {
	// 	models.NewErrorResponse(c, http.StatusBadRequest, "empty param standardId", "стандарт не задан")
	// 	return
	// }
	constructionId := c.Query("constructionId")
	if constructionId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param constructionId", "конструкция не задана")
		return
	}
	// fillerId := c.Query("fillerId")
	configuration := c.Query("configuration")
	// changeStandard := c.Query("changeStandard")

	data, err := h.putgApi.GetData(c, &putg_api.GetPutgData{
		StandardId:     standardId,
		ConstructionId: constructionId,
		Configuration:  configuration,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: data})
}
