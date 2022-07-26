package pro

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/pro_model"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func (h *Handler) initSizeIntRoutes(api *gin.RouterGroup) {
	size := api.Group("/size-interview")
	{
		size.GET("/", h.getSizesInt)
		size.GET("/all", h.getAllSizesInt)
		size = size.Group("/", h.middleware.AccessForProAdmin)
		{
			size.POST("/", h.createSizeInt)
			size.POST("/file", h.createSizeIntFromFile)
			size.PUT("/:id", h.updateSizeInt)
			size.DELETE("/:id", h.deleteSizeInt)
			size.DELETE("/all", h.deleteAllSizeInt)
		}
	}
}

// @Summary Get Sizes Int
// @Tags Sealur Pro -> sizes-interview
// @Description получение размеров (для опроса)
// @ModuleID getSizesInt
// @Accept json
// @Produce json
// @Param flange query string true "flange"
// @Param typeFlId query string true "flange type id"
// @Param row query string true "row"
// @Success 200 {object} models.DataResponse{data=pro_api.SizeIntResponse}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes-interview [get]
func (h *Handler) getSizesInt(c *gin.Context) {
	flange := c.Query("flange")
	if flange == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty flange", "empty flange param")
		return
	}
	typeFlId := c.Query("typeFlId")
	if typeFlId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty flange type id", "empty flange type id param")
		return
	}
	rowStr := c.Query("row")
	if rowStr == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty row", "empty row")
		return
	}
	row, err := strconv.Atoi(rowStr)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty row", "empty row")
		return
	}

	sizes, err := h.proClient.GetSizeInt(c, &pro_api.GetSizesIntRequest{
		FlangeId: flange,
		TypeFl:   typeFlId,
		Row:      int32(row),
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: sizes, Count: len(sizes.Sizes)})
}

// @Summary Get All Sizes Int
// @Tags Sealur Pro -> sizes-interview
// @Description получение размеров (для опроса)
// @ModuleID getAllSizesInt
// @Accept json
// @Produce json
// @Param flange query string true "flange"
// @Param typeFlId query string true "flange type id"
// @Success 200 {object} models.DataResponse{data=pro_api.SizeIntResponse}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes-interview/all [get]
func (h *Handler) getAllSizesInt(c *gin.Context) {
	flange := c.Query("flange")
	if flange == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty flange", "empty flange param")
		return
	}
	typeFlId := c.Query("typeFlId")
	if typeFlId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty flange type id", "empty flange type id param")
		return
	}

	sizes, err := h.proClient.GetAllSizeInt(c, &pro_api.GetAllSizeIntRequest{
		FlangeId: flange,
		TypeFl:   typeFlId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: sizes, Count: len(sizes.Sizes)})
}

// @Summary Create Size Int
// @Tags Sealur Pro -> sizes-interview
// @Security ApiKeyAuth
// @Description создание размеров (для опроса)
// @ModuleID createSizeInt
// @Accept json
// @Produce json
// @Param data body pro_model.SizeIntDTO true "size int info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes-interview [post]
func (h *Handler) createSizeInt(c *gin.Context) {
	var dto pro_model.SizeIntDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	size, err := h.proClient.CreateSizeInt(c, &pro_api.CreateSizeIntRequest{
		FlangeId:  dto.Flange,
		TypeFl:    dto.TypeFlId,
		Dy:        dto.Dy,
		Py:        dto.Py,
		D1:        dto.D1,
		D2:        dto.D2,
		DUp:       dto.DUp,
		D:         dto.D,
		H1:        dto.H1,
		H2:        dto.H2,
		Bolt:      dto.Bolt,
		CountBolt: dto.CountBolt,
		Number:    dto.Count,
		Row:       dto.Row,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/sizes-interview/%s", size.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: size.Id, Message: "Created"})
}

// @Summary Create Size Int From File
// @Tags Sealur Pro -> sizes-interview
// @Security ApiKeyAuth
// @Description создание размеров (для опроса) из файла excel
// @ModuleID createSizeFromFile
// @Accept multipart/form-data
// @Produce json
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes/file [post]
func (h *Handler) createSizeIntFromFile(c *gin.Context) {
	fileHeader, err := c.FormFile("sizes")
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "error while opening file")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "error while opening file")
		return
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "error while reading file")
		return
	}

	sheetName := f.GetSheetName(f.GetActiveSheetIndex())

	req := make([]*pro_api.CreateSizeIntRequest, 0, 50)

	rows, err := f.Rows(sheetName)
	if err != nil {
		logger.Error(err)
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "error while reading file")
		return
	}
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			logger.Error(err)
		}

		if len(row) == 0 || row[0] == "" || row[0] == "count" {
			continue
		}

		count, err := strconv.Atoi(row[0])
		if err != nil {
			logger.Error("count empty")
			count = 0
		}
		countBolt, err := strconv.Atoi(row[12])
		if err != nil {
			logger.Error("count empty")
			countBolt = 0
		}
		rowCount, err := strconv.Atoi(row[13])
		if err != nil {
			logger.Error("count empty")
			rowCount = 0
		}

		req = append(req, &pro_api.CreateSizeIntRequest{
			Number:    int32(count),
			FlangeId:  row[1],
			Dy:        row[2],
			Py:        row[3],
			TypeFl:    row[4],
			DUp:       row[5],
			D1:        row[6],
			D2:        row[7],
			D:         row[8],
			H1:        row[9],
			H2:        row[10],
			Bolt:      row[11],
			CountBolt: int32(countBolt),
			Row:       int32(rowCount),
		})
	}
	if err = rows.Close(); err != nil {
		logger.Error(err)
	}

	_, err = h.proClient.CreateManySizesInt(c, &pro_api.CreateSizesIntRequest{Sizes: req})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

// @Summary Update Size Int
// @Tags Sealur Pro -> sizes-interview
// @Security ApiKeyAuth
// @Description обновление размеров (для опроса)
// @ModuleID updateSizeInt
// @Accept json
// @Produce json
// @Param data body pro_model.SizeIntDTO true "size int info"
// @Param id path string true "size int id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes-interview/{id} [put]
func (h *Handler) updateSizeInt(c *gin.Context) {
	var dto pro_model.SizeIntDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	size, err := h.proClient.UpdateSizeInt(c, &pro_api.UpdateSizeIntRequest{
		Id:        id,
		FlangeId:  dto.Flange,
		TypeFl:    dto.TypeFlId,
		Dy:        dto.Dy,
		Py:        dto.Py,
		D1:        dto.D1,
		D2:        dto.D2,
		DUp:       dto.DUp,
		D:         dto.D,
		H1:        dto.H1,
		H2:        dto.H2,
		Bolt:      dto.Bolt,
		CountBolt: dto.CountBolt,
		Row:       dto.Row,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: size.Id, Message: "Updated"})
}

// @Summary Delete Size Int
// @Tags Sealur Pro -> sizes-interview
// @Security ApiKeyAuth
// @Description удаление размеров (для опроса)
// @ModuleID deleteSizeInt
// @Accept json
// @Produce json
// @Param id path string true "size id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes-interview/{id} [delete]
func (h *Handler) deleteSizeInt(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	size, err := h.proClient.DeleteSizeInt(c, &pro_api.DeleteSizeIntRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: size.Id, Message: "Deleted"})
}

// @Summary Delete All Size Int
// @Tags Sealur Pro -> sizes-interview
// @Security ApiKeyAuth
// @Description удаление размеров (для опроса)
// @ModuleID deleteAllSizeInt
// @Accept json
// @Produce json
// @Param flange query string true "flange"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes-interview/all [delete]
func (h *Handler) deleteAllSizeInt(c *gin.Context) {
	flange := c.Query("flange")
	if flange == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty flange", "empty flange param")
		return
	}

	_, err := h.proClient.DeleteAllSizeInt(c, &pro_api.DeleteAllSizeIntRequest{FlangeId: flange})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Deleted"})
}
