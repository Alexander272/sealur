package file

import (
	"bufio"
	"bytes"
	"io"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto/proto_file"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) initFilesRoutes(api *gin.RouterGroup) {
	drawing := api.Group("/drawings")
	{
		drawing.POST("/:backet", h.createDrawing)
		drawing.GET("/:backet/:group/:id/:name", h.getDrawing)
		drawing.DELETE("/:backet/:group/:id/:name", h.deleteDrawing)
	}
}

// @Summary Get Drawing
// @Tags Files -> drawing
// @Security ApiKeyAuth
// @Description создание чертежа
// @ModuleID getDrawing
// @Accept json
// @Produce multipart/form-data
// @Param name path string true "drawing name"
// @Param id path string true "drawing id"
// @Param group path string true "drawing group"
// @Param backet path string true "backet"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /files/drawings/{backet}/{group}/{id}/{name} [get]
func (h *Handler) getDrawing(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty name", "empty name param")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	group := c.Param("group")
	if group == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty group", "empty group param")
		return
	}

	backet := c.Param("backet")
	if backet == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty backet", "empty backet param")
		return
	}

	stream, err := h.fileClient.Download(c, &proto_file.FileDownloadRequest{
		Id:     id,
		Name:   name,
		Backet: backet,
		Group:  group,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	req, err := stream.Recv()
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	meta := req.GetMetadata()
	imageData := bytes.Buffer{}

	logger.Debug(meta)

	for {
		logger.Debug("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			logger.Debug("no more data")
			break
		}

		if err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
			return
		}

		chunk := req.GetFile().Content

		_, err = imageData.Write(chunk)
		if err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
			return
		}
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+name)
	c.Data(http.StatusOK, meta.Type, imageData.Bytes())
}

// @Summary Create Drawing
// @Tags Files -> drawing
// @Security ApiKeyAuth
// @Description создание чертежа
// @ModuleID createDrawing
// @Accept multipart/form-data
// @Produce json
// @Param backet path string true "backet"
// @Param group body string false "group image"
// @Success 201 {object} models.FileResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /files/drawings/{backet} [post]
func (h *Handler) createDrawing(c *gin.Context) {
	backet := c.Param("backet")
	if backet == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty backet", "empty backet param")
		return
	}

	file, err := c.FormFile("drawing")
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "error getting file")
		return
	}

	f, err := file.Open()
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "file read error")
		return
	}
	defer f.Close()

	group := c.Request.FormValue("group")
	if group == "" {
		gUuid, err := uuid.NewUUID()
		if err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "failed to generate group id")
			return
		}
		group = gUuid.String()
	}

	fileType := file.Header.Get("Content-Type")

	reqMeta := &proto_file.FileUploadRequest{
		Request: &proto_file.FileUploadRequest_Metadata{
			Metadata: &proto_file.MetaData{
				Name:   file.Filename,
				Type:   fileType,
				Size:   file.Size,
				Group:  group,
				Backet: backet,
			},
		},
	}

	stream, err := h.fileClient.Upload(c)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "error while connect wuth service")
		return
	}

	err = stream.Send(reqMeta)
	if err != nil {
		logger.Errorf("cannot send image info to server: %w %w", err, stream.RecvMsg(nil))
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "cannot send image info to server")
		return
	}

	reader := bufio.NewReader(f)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Errorf("cannot read chunk to buffer: %w", err)
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "cannot send image to server")
			return
		}

		reqChunk := &proto_file.FileUploadRequest{
			Request: &proto_file.FileUploadRequest_File{
				File: &proto_file.File{
					Content: buffer[:n],
				},
			},
		}

		err = stream.Send(reqChunk)
		if err != nil {
			logger.Errorf("cannot send chunk to server: %w", err)
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "cannot send image to server")
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		logger.Errorf("cannot receive response: %w", err)
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "cannot receive response")
		return
	}

	c.JSON(http.StatusCreated, models.FileResponse{
		Id:       res.Id,
		Name:     res.Name,
		OrigName: res.OrigName,
		Link:     res.Url,
		Group:    group,
	})
}

// @Summary Delete Drawing
// @Tags Files -> drawing
// @Security ApiKeyAuth
// @Description удаление чертежа
// @ModuleID deleteDrawing
// @Accept json
// @Produce json
// @Param name path string true "drawing name"
// @Param id path string true "drawing id"
// @Param group path string true "drawing group"
// @Param backet path string true "backet"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /files/drawings/{backet}/{group}/{id}/{name} [delete]
func (h *Handler) deleteDrawing(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty name", "empty name param")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	group := c.Param("group")
	if group == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty group", "empty group param")
		return
	}

	backet := c.Param("backet")
	if backet == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty backet", "empty backet param")
		return
	}

	_, err := h.fileClient.Delete(c, &proto_file.FileDeleteRequest{
		Id:     id,
		Name:   name,
		Backet: backet,
		Group:  group,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Deleted"})
}
