package http

import (
	"context"
	"net/http"

	"github.com/Alexander272/sealur/file_service/internal/config"
	"github.com/Alexander272/sealur/file_service/internal/models"
	"github.com/Alexander272/sealur/file_service/internal/service"
	"github.com/Alexander272/sealur/file_service/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init(conf *config.Config) *gin.Engine {
	router := gin.Default()

	files := router.Group("/files")
	{
		// files.GET("/:name", h.getFile)
		files.GET("/:backet/:group/:id/:name", h.getFile)
	}

	return router
}

func (h *Handler) getFile(c *gin.Context) {
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
	if group == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty backet", "empty backet param")
		return
	}

	logger.Debug(name, " ", id, " ", group, " ", backet)

	file, err := h.services.GetFile(context.Background(), backet, group, id, name)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "error getting file")
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+name)
	c.Data(http.StatusOK, file.ContentType, file.Bytes)
}

// func (h *Handler) getFile(c *gin.Context) {
// 	name := c.Param("name")
// 	if name == "" {
// 		models.NewErrorResponse(c, http.StatusBadRequest, "empty name", "empty name param")
// 		return
// 	}

// 	id := c.Query("id")
// 	if id == "" {
// 		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
// 		return
// 	}

// 	group := c.Query("group")
// 	if group == "" {
// 		models.NewErrorResponse(c, http.StatusBadRequest, "empty group", "empty group param")
// 		return
// 	}

// 	backet := c.Query("backet")
// 	if group == "" {
// 		models.NewErrorResponse(c, http.StatusBadRequest, "empty backet", "empty backet param")
// 		return
// 	}

// 	logger.Debug(name, " ", id, " ", group, " ", backet)

// 	file, err := h.services.GetFile(context.Background(), backet, group, id, name)
// 	if err != nil {
// 		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "error getting file")
// 		return
// 	}

// 	c.Header("Content-Description", "File Transfer")
// 	c.Header("Content-Transfer-Encoding", "binary")
// 	c.Header("Content-Disposition", "attachment; filename="+name)
// 	c.Data(http.StatusOK, file.ContentType, file.Bytes)
// }
