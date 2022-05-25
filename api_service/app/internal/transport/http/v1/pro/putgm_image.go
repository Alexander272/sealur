package pro

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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
// @Accept multipart/form-data
// @Accept json
// @Produce json
// @Param data body models.PutgmImageDTO true "putgm image info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putgm-image [post]
func (h *Handler) createPutgmImage(c *gin.Context) {
	var dto models.PutgmImageDTO

	if c.Request.FormValue("gasket") != "" {
		dto.Gasket = c.Request.FormValue("gasket")
		dto.Form = c.Request.FormValue("form")
		dto.Url = c.Request.FormValue("url")
	} else {
		if err := c.BindJSON(&dto); err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
			return
		}
	}

	var folder string
	if dto.Form == "Round" {
		folder = "putgm/construction"
	} else {
		folder = "putgm/half"
	}

	err := h.saveFile(c, folder)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "error while saving file")
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

func (h *Handler) saveFile(c *gin.Context, folder string) error {
	file, err := c.FormFile("drawing")
	if err != nil {
		return err
	}

	if err := c.SaveUploadedFile(file, filepath.Join("images", folder, file.Filename)); err != nil {
		return err
	}

	return nil
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
// @Param file query string true "file name"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/putgm-image/{id} [delete]
func (h *Handler) deletePutgmImage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	file := c.Query("file")
	if file == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty file", "empty file param")
		return
	}
	os.Remove(filepath.Join("images", file))

	image, err := h.proClient.DeletePutgmImage(c, &proto.DeletePutgmImageRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: image.Id, Message: "Deleted"})
}
