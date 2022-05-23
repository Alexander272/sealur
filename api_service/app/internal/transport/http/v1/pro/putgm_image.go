package pro

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initPutgmImageRoutes(api *gin.RouterGroup) {
	putgmImage := api.Group("putgm-image")
	{
		putgmImage.GET("/", h.getPutgmImage)
		putgmImage.POST("/", h.createPutgmImage)
		putgmImage.PUT("/:id", h.updatePutgmImage)
		putgmImage.DELETE("/:id", h.deletePutgmImage)
	}
}

// @Summary Get Putgm Image
// @Tags Sealur Pro -> putgm-image
// @Security ApiKeyAuth
// @Description получение списка чертежей для путгм
// @ModuleID getPutgmImage
// @Accept json
// @Produce json
// @Param form query string true "form"
// @Success 200 {object} models.DataResponse{data=[]proto.PutgmImage}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putgm-image [get]
func (h *Handler) getPutgmImage(c *gin.Context) {
	form := c.Query("form")
	if form == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty form", "empty form param")
		return
	}

	images, err := h.proClient.GetPutgmImage(c, &proto.GetPutgmImageRequest{Form: form})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: images.PutgmImage, Count: len(images.PutgmImage)})
}

// @Summary Create Putgm Image
// @Tags Sealur Pro -> putgm-image
// @Security ApiKeyAuth
// @Description создание чертежа путгм
// @ModuleID createPutgmImage
// @Accept json
// @Produce json
// @Param data body models.PutgmImageDTO true "putgm image info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putgm-image [post]
func (h *Handler) createPutgmImage(c *gin.Context) {
	//TODO дописать добавление картинки

	var dto models.PutgmImageDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	request := &proto.CreatePutgmImageRequest{
		Form:   dto.Form,
		Gasket: dto.Gasket,
		Url:    dto.Url,
	}

	image, err := h.proClient.CreatePutgmImage(c, request)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/putgm-image/%s", image.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: image.Id, Message: "Created"})
}

// @Summary Update Putgm Image
// @Tags Sealur Pro -> putgm-image
// @Security ApiKeyAuth
// @Description обновление чертежа путгм
// @ModuleID updatePutgmImage
// @Accept json
// @Produce json
// @Param id path string true "putgm image id"
// @Param data body models.PutgmImageDTO true "putgm image info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putgm-image/{id} [put]
func (h *Handler) updatePutgmImage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto models.PutgmImageDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	request := &proto.UpdatePutgmImageRequest{
		Id:     id,
		Form:   dto.Form,
		Gasket: dto.Gasket,
		Url:    dto.Url,
	}

	image, err := h.proClient.UpdatePutgmImage(c, request)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: image.Id, Message: "Updated"})
}

// @Summary Delete Putgm Image
// @Tags Sealur Pro -> putgm-image
// @Security ApiKeyAuth
// @Description удаление чертежа путг
// @ModuleID deletePutgmImage
// @Accept json
// @Produce json
// @Param id path string true "putgm image id"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putgm-image/{id} [delete]
func (h *Handler) deletePutgmImage(c *gin.Context) {
	//TODO дописать удаление картинки
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	image, err := h.proClient.DeletePutgmImage(c, &proto.DeletePutgmImageRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: image.Id, Message: "Deleted"})
}
