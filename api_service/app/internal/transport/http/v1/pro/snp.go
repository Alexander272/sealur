package pro

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initSNPRoutes(api *gin.RouterGroup) {
	snp := api.Group("/snp")
	{
		snp.GET("/", h.getSNP)
		snp.POST("/", h.createSNP)
		snp.PUT("/:id", h.updateSNP)
		snp.DELETE("/:id", h.deleteSNP)
	}
}

// @Summary Get SNP
// @Tags Sealur Pro -> snp
// @Security ApiKeyAuth
// @Description получение прокладок снп
// @ModuleID getSNP
// @Accept json
// @Produce json
// @Param standId query string true "stand id"
// @Param flangeId query string true "flange id"
// @Success 200 {object} models.DataResponse{data=[]proto.SNP}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/snp [get]
func (h *Handler) getSNP(c *gin.Context) {
	standId := c.Query("standId")
	if standId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty standId param")
		return
	}
	flangeId := c.Query("flangeId")
	if flangeId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty flangeId param")
		return
	}

	snp, err := h.proClient.GetSNP(c, &proto.GetSNPRequest{
		StandId:  standId,
		FlangeId: flangeId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: snp.Snp, Count: len(snp.Snp)})
}

// @Summary Create SNP
// @Tags Sealur Pro -> snp
// @Security ApiKeyAuth
// @Description создание прокладки снп
// @ModuleID createSNP
// @Accept json
// @Produce json
// @Param data body models.SNPDTO true "snp info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/snp [post]
func (h *Handler) createSNP(c *gin.Context) {
	var dto models.SNPDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	snp, err := h.proClient.CreateSNP(c, &proto.CreateSNPRequest{
		StandId:   dto.StandId,
		FlangeId:  dto.FlangeId,
		TypeFlId:  dto.TypeFlId,
		TypePr:    dto.TypePr,
		Fillers:   dto.Fillers,
		Materials: dto.Materials,
		DefMat:    dto.DefMat,
		Mounting:  dto.Mounting,
		Graphite:  dto.Graphite,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/snp/%s", snp.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: snp.Id, Message: "Created"})
}

// @Summary Update SNP
// @Tags Sealur Pro -> snp
// @Security ApiKeyAuth
// @Description обновление прокладки снп
// @ModuleID updateSNP
// @Accept json
// @Produce json
// @Param data body models.SNPDTO true "snp info"
// @Param id path string true "snp id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/snp/{id} [put]
func (h *Handler) updateSNP(c *gin.Context) {
	var dto models.SNPDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	snp, err := h.proClient.UpdateSNP(c, &proto.UpdateSNPRequest{
		Id:        id,
		StandId:   dto.StandId,
		FlangeId:  dto.FlangeId,
		TypeFlId:  dto.TypeFlId,
		TypePr:    dto.TypePr,
		Fillers:   dto.Fillers,
		Materials: dto.Materials,
		DefMat:    dto.DefMat,
		Mounting:  dto.Mounting,
		Graphite:  dto.Graphite,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: snp.Id, Message: "Updated"})
}

// @Summary Delete SNP
// @Tags Sealur Pro -> snp
// @Security ApiKeyAuth
// @Description удаление прокладки снп
// @ModuleID deleteSNP
// @Accept json
// @Produce json
// @Param id path string true "snp id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/snp/{id} [delete]
func (h *Handler) deleteSNP(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	snp, err := h.proClient.DeleteSNP(c, &proto.DeleteSNPRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: snp.Id, Message: "Deleted"})
}
