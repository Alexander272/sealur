package pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/pro_model"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initInterviewRoutes(api *gin.RouterGroup) {
	materials := api.Group("/interview")
	{
		materials.POST("/", h.sendInterview)
	}
}

// @Summary Send Interview
// @Tags Sealur Pro -> interview
// @Description отправление данных опроса
// @ModuleID sendInterview
// @Accept json
// @Produce json
// @Param data body pro_model.Interview true "interview info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/interview [post]
func (h *Handler) sendInterview(c *gin.Context) {
	var dto pro_model.Interview
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	var drawing *pro_api.Drawing = nil
	if dto.Drawing != nil {
		drawing = &pro_api.Drawing{
			Id:       dto.Drawing.Id,
			Name:     dto.Drawing.Name,
			OrigName: dto.Drawing.OrigName,
			Group:    dto.Drawing.Group,
			Link:     dto.Drawing.Link,
		}
	}

	size := &pro_api.SizesInt{
		Dy:        dto.Sizes.Dy,
		Py:        dto.Sizes.Py,
		DUp:       dto.Sizes.DUp,
		D1:        dto.Sizes.D1,
		D2:        dto.Sizes.D2,
		D:         dto.Sizes.D,
		H1:        dto.Sizes.H1,
		H2:        dto.Sizes.H2,
		Bolt:      dto.Sizes.Bolt,
		CountBolt: dto.Sizes.CountBolt,
		DIn:       dto.Sizes.DIn,
		DOut:      dto.Sizes.DOut,
		H:         dto.Sizes.H,
	}

	_, err := h.proClient.SendInterview(c, &pro_api.SendInterviewRequest{
		Organization:  dto.Organization,
		Name:          dto.Name,
		Email:         dto.Email,
		City:          dto.City,
		Position:      dto.Position,
		Phone:         dto.Phone,
		Techprocess:   dto.Techprocess,
		Equipment:     dto.Equipment,
		Seal:          dto.Seal,
		Consumer:      dto.Consumer,
		Factory:       dto.Factory,
		Developer:     dto.Developer,
		Flange:        dto.Flange,
		TypeFl:        dto.TypeFl,
		Type:          dto.Type,
		DiffFrom:      dto.DiffFrom,
		DiffTo:        dto.DiffTo,
		PresWork:      dto.PresWork,
		PresTest:      dto.PresTest,
		Pressure:      dto.Pressure,
		Environ:       dto.Environ,
		TempWorkPipe:  dto.TempWorkPipe,
		PresWorkPipe:  dto.PresWorkPipe,
		EnvironPipe:   dto.EnvironPipe,
		TempWorkAnn:   dto.TempWorkAnn,
		PresWorkAnn:   dto.PresWorkAnn,
		EnvironAnn:    dto.EnvironAnn,
		Material:      dto.Material,
		BoltMaterial:  dto.BoltMaterial,
		Lubricant:     dto.Lubricant,
		Along:         dto.Along,
		Across:        dto.Across,
		NonFlatness:   dto.NonFlatness,
		Mounting:      dto.Mounting,
		Condition:     dto.Condition,
		Period:        dto.Period,
		Abrasive:      dto.Abrasive,
		Crystallized:  dto.Crystallized,
		Penetrating:   dto.Penetrating,
		DrawingNumber: dto.DrawingNumber,
		Info:          dto.Info,
		Drawing:       drawing,
		Sizes:         size,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Sent successfully"})
}
