package moment

import (
	"net/http"
	"strconv"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initBoltsRoutes(api *gin.RouterGroup) {
	bolts := api.Group("/bolts", h.middleware.UserIdentity)
	{
		bolts.GET("/", h.getBolts)
		bolts.GET("/all", h.getAllBolts)
		bolts = bolts.Group("/", h.middleware.AccessForMomentAdmin)
		{
			bolts.POST("/", h.createBolt)
			bolts.POST("/several", h.createBolts)
			bolts.PUT("/:id", h.updateBolt)
			bolts.DELETE("/:id", h.deleteBolt)
		}
	}
}

// @Summary Get Bolts
// @Tags Sealur Moment -> bolts
// @Security ApiKeyAuth
// @Description получение болтов
// @ModuleID getBolts
// @Accept json
// @Produce json
// @Param isInch query bool false "is inch"
// @Success 200 {object} models.DataResponse{Data=[]flange_model.Bolt}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/bolts/ [get]
func (h *Handler) getBolts(c *gin.Context) {
	defIsInch := false
	isInchQeury := c.Query("isInch")
	if isInchQeury != "" {
		isInch, err := strconv.ParseBool(isInchQeury)
		if err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
			return
		}
		defIsInch = isInch
	}

	bolts, err := h.flangeClient.GetBolts(c, &flange_api.GetBoltsRequest{IsInch: defIsInch})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: bolts.Bolts, Count: len(bolts.Bolts)})
}

// @Summary Get All Bolts
// @Tags Sealur Moment -> bolts
// @Security ApiKeyAuth
// @Description получение болтов
// @ModuleID getAllBolts
// @Accept json
// @Produce json
// @Param isInch query bool false "is inch"
// @Success 200 {object} models.DataResponse{Data=[]flange_model.Bolt}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/bolts/all [get]
func (h *Handler) getAllBolts(c *gin.Context) {
	bolts, err := h.flangeClient.GetAllBolts(c, &flange_api.GetBoltsRequest{})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: bolts.Bolts, Count: len(bolts.Bolts)})
}

// @Summary Create Bolt
// @Tags Sealur Moment -> bolts
// @Security ApiKeyAuth
// @Description создание болта
// @ModuleID createBolt
// @Accept json
// @Produce json
// @Param bolt body moment_model.BoltDTO true "bolt info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/bolts/ [post]
func (h *Handler) createBolt(c *gin.Context) {
	var dto moment_model.BoltDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.flangeClient.CreateBolt(c, &flange_api.CreateBoltRequest{
		Title:    dto.Title,
		Diameter: dto.Diameter,
		Area:     dto.Area,
		IsInch:   dto.IsInch,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

// @Summary Create Bolts
// @Tags Sealur Moment -> bolts
// @Security ApiKeyAuth
// @Description создание болтов
// @ModuleID createBolts
// @Accept json
// @Produce json
// @Param bolt body []moment_model.BoltDTO true "bolt info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/bolts/several [post]
func (h *Handler) createBolts(c *gin.Context) {
	var dto []moment_model.BoltDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	bolts := []*flange_api.CreateBoltRequest{}
	for _, bd := range dto {
		bolts = append(bolts, &flange_api.CreateBoltRequest{
			Title:    bd.Title,
			Diameter: bd.Diameter,
			Area:     bd.Area,
			IsInch:   bd.IsInch,
		})
	}

	_, err := h.flangeClient.CreateBolts(c, &flange_api.CreateBoltsRequest{
		Bolts: bolts,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

// @Summary Update Bolt
// @Tags Sealur Moment -> bolts
// @Security ApiKeyAuth
// @Description обновление среды
// @ModuleID updateBolt
// @Accept json
// @Produce json
// @Param id path string true "bolt id"
// @Param bolt body moment_model.BoltDTO true "bolt info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/bolts/{id} [put]
func (h *Handler) updateBolt(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.BoltDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.flangeClient.UpdateBolt(c, &flange_api.UpdateBoltRequest{
		Id:       id,
		Title:    dto.Title,
		Diameter: dto.Diameter,
		Area:     dto.Area,
		IsInch:   dto.IsInch,
	})
	if err != nil {
		models.NewErrorResponseWithCode(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

// @Summary Delete Bolt
// @Tags Sealur Moment -> bolts
// @Security ApiKeyAuth
// @Description Удаление болта
// @ModuleID deleteBolt
// @Accept json
// @Produce json
// @Param id path string true "bolt id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/bolts/{id} [delete]
func (h *Handler) deleteBolt(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.flangeClient.DeleteBolt(c, &flange_api.DeleteBoltRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
