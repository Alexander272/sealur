package file

import (
	"bufio"
	"io"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto/proto_file"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initFilesRoutes(api *gin.RouterGroup) {

	// order.GET("/", h.getPutg)
	// order.POST("/", h.createPutg)
	// order.PUT("/:id", h.updatePutg)
	// order.DELETE("/:id", h.deletePutg)

	drawing := api.Group("/drawings")
	{
		drawing.POST("/", h.createDrawing)
	}
}

// @Summary Create Drawing
// @Tags Files -> drawing
// @Security ApiKeyAuth
// @Description создание чертежа
// @ModuleID createDrawing
// @Accept multipart/form-data
// @Produce json
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /files/drawing [post]
func (h *Handler) createDrawing(c *gin.Context) {
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

	testUUID := "a053dc2e-dfd3-11ec-9d64-0242ac120002"

	fileType := file.Header.Get("Content-Type")

	reqMeta := &proto_file.FileUploadRequest{
		Request: &proto_file.FileUploadRequest_Metadata{
			Metadata: &proto_file.MetaData{
				Name:   file.Filename,
				Type:   fileType,
				Size:   file.Size,
				Uuid:   testUUID,
				Backet: "pro",
			},
		},
	}
	// reqFile := &proto_file.FileUploadRequest{
	// 	Request: &proto_file.FileUploadRequest_File{},
	// }

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

	logger.Debug(res)
	logger.Debug(*res)

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}
