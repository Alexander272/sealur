package pro

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initPutgImageRoutes(api *gin.RouterGroup) {
	st := api.Group("putg-image")
	{
		st.GET("/", h.getPutgImage)
		st.POST("/", h.createPutgImage)
		st.PUT("/:id", h.updatePutgImage)
		st.DELETE("/:id", h.deletePutgImage)
	}
}

// @Summary Get Putg Image
// @Tags Sealur Pro -> putg-image
// @Security ApiKeyAuth
// @Description получение списка чертежей для путг
// @ModuleID getPutgImage
// @Accept json
// @Produce json
// @Param form query string true "form"
// @Success 200 {object} models.DataResponse{data=[]proto.PutgImage}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putg-image [get]
func (h *Handler) getPutgImage(c *gin.Context) {
	form := c.Query("form")
	if form == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty form", "empty form param")
		return
	}

	images, err := h.proClient.GetPutgImage(c, &proto.GetPutgImageRequest{Form: form})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: images.PutgImage, Count: len(images.PutgImage)})
}

// @Summary Create Putg Image
// @Tags Sealur Pro -> putg-image
// @Security ApiKeyAuth
// @Description создание чертежа путг
// @ModuleID createPutgImage
// @Accept json
// @Produce json
// @Param data body models.PutgImageDTO true "putg image info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putg-image [post]
func (h *Handler) createPutgImage(c *gin.Context) {
	//TODO дописать добавление картинки

	var dto models.PutgImageDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	request := &proto.CreatePutgImageRequest{
		Form:   dto.Form,
		Gasket: dto.Gasket,
		Url:    dto.Url,
	}

	image, err := h.proClient.CreatePutgImage(c, request)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/putg-image/%s", image.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: image.Id, Message: "Created"})
}

// @Summary Update Putg Image
// @Tags Sealur Pro -> putg-image
// @Security ApiKeyAuth
// @Description обновление чертежа путг
// @ModuleID updatePutgImage
// @Accept json
// @Produce json
// @Param id path string true "putg image id"
// @Param data body models.PutgImageDTO true "putg image info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putg-image/{id} [put]
func (h *Handler) updatePutgImage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto models.PutgImageDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	request := &proto.UpdatePutgImageRequest{
		Id:     id,
		Form:   dto.Form,
		Gasket: dto.Gasket,
		Url:    dto.Url,
	}

	image, err := h.proClient.UpdatePutgImage(c, request)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: image.Id, Message: "Updated"})
}

// @Summary Delete Putg Image
// @Tags Sealur Pro -> putg-image
// @Security ApiKeyAuth
// @Description удаление чертежа путг
// @ModuleID deletePutgImage
// @Accept json
// @Produce json
// @Param id path string true "st/fl id"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putg-image/{id} [delete]
func (h *Handler) deletePutgImage(c *gin.Context) {
	//TODO дописать удаление картинки
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	image, err := h.proClient.DeletePutgImage(c, &proto.DeletePutgImageRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: image.Id, Message: "Deleted"})
}
