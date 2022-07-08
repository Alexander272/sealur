package pro

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func (h *Handler) initSizeRoutes(api *gin.RouterGroup) {
	sizes := api.Group("/sizes", h.middleware.UserIdentity)
	{
		sizes.GET("/", h.getSizes)
		sizes.GET("/all", h.getAllSizes)
		sizes = sizes.Group("/", h.middleware.AccessForProAdmin)
		{
			sizes.POST("/", h.createSize)
			sizes.POST("/file", h.createSizeFromFile)
			sizes.PUT("/:id", h.updateSize)
			sizes.DELETE("/:id", h.deleteSize)
			sizes.DELETE("/all", h.deleteAllSize)
		}
	}
}

// @Summary Get Sizes
// @Tags Sealur Pro -> sizes
// @Description получение размеров
// @ModuleID getSizes
// @Accept json
// @Produce json
// @Param flange query string true "flange"
// @Param typeFlId query string true "flange type id"
// @Param standId query string true "standarts id"
// @Param typePr query string true "type"
// @Success 200 {object} models.DataResponse{data=proto.SizeResponse}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes [get]
func (h *Handler) getSizes(c *gin.Context) {
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
	standId := c.Query("standId")
	if standId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty stand id", "empty standarts id param")
		return
	}
	typePr := c.Query("typePr")
	if typePr == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty type", "empty type lining")
		return
	}

	sizes, err := h.proClient.GetSizes(c, &proto.GetSizesRequest{
		Flange:   flange,
		TypeFlId: typeFlId,
		TypePr:   typePr,
		StandId:  standId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: sizes, Count: len(sizes.Sizes)})
}

// @Summary Get All Sizes
// @Tags Sealur Pro -> sizes
// @Description получение размеров
// @ModuleID getSizes
// @Accept json
// @Produce json
// @Param flange query string true "flange"
// @Param typeFlId query string true "flange type id"
// @Param standId query string true "standarts id"
// @Param typePr query string true "type"
// @Success 200 {object} models.DataResponse{data=proto.SizeResponse}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes/all [get]
func (h *Handler) getAllSizes(c *gin.Context) {
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
	standId := c.Query("standId")
	if standId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty stand id", "empty standarts id param")
		return
	}
	typePr := c.Query("typePr")
	if typePr == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty type", "empty type lining")
		return
	}

	sizes, err := h.proClient.GetAllSizes(c, &proto.GetSizesRequest{
		Flange:   flange,
		TypeFlId: typeFlId,
		TypePr:   typePr,
		StandId:  standId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: sizes, Count: len(sizes.Sizes)})
}

// @Summary Create Size
// @Tags Sealur Pro -> sizes
// @Security ApiKeyAuth
// @Description создание размеров
// @ModuleID createSize
// @Accept json
// @Produce json
// @Param data body models.SizesDTO true "size info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes [post]
func (h *Handler) createSize(c *gin.Context) {
	var dto models.SizesDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	size, err := h.proClient.CreateSize(c, &proto.CreateSizeRequest{
		Flange:   dto.Flange,
		TypeFlId: dto.TypeFlId,
		Dn:       dto.Dn,
		Pn:       dto.Pn,
		TypePr:   dto.TypePr,
		StandId:  dto.StandId,
		D4:       dto.D4,
		D3:       dto.D3,
		D2:       dto.D2,
		D1:       dto.D1,
		H:        dto.H,
		S2:       dto.S2,
		S3:       dto.S3,
		Number:   dto.Number,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/sizes/%s", size.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: size.Id, Message: "Created"})
}

// @Summary Create Size From File
// @Tags Sealur Pro -> sizes
// @Security ApiKeyAuth
// @Description создание размеров из файла excel
// @ModuleID createSizeFromFile
// @Accept multipart/form-data
// @Produce json
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes/file [post]
func (h *Handler) createSizeFromFile(c *gin.Context) {
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

	req := make([]*proto.CreateSizeRequest, 0, 50)

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

		var s2, s3 string
		if len(row) <= 12 {
			s2, s3 = "", ""
		} else {
			s2 = row[12]
			s3 = row[13]
		}

		req = append(req, &proto.CreateSizeRequest{
			Number:   int32(count),
			Flange:   row[1],
			Dn:       row[2],
			Pn:       row[3],
			TypePr:   row[4],
			StandId:  row[5],
			TypeFlId: row[6],
			D4:       row[7],
			D3:       row[8],
			D2:       row[9],
			D1:       row[10],
			H:        row[11],
			S2:       s2,
			S3:       s3,
		})
	}
	if err = rows.Close(); err != nil {
		logger.Error(err)
	}

	_, err = h.proClient.CreateManySizes(c, &proto.CreateSizesRequest{Sizes: req})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

// @Summary Update Size
// @Tags Sealur Pro -> sizes
// @Security ApiKeyAuth
// @Description обновление размеров
// @ModuleID updateSize
// @Accept json
// @Produce json
// @Param data body models.SizesDTO true "size info"
// @Param id path string true "size id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes/{id} [put]
func (h *Handler) updateSize(c *gin.Context) {
	var dto models.SizesDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	size, err := h.proClient.UpdateSize(c, &proto.UpdateSizeRequest{
		Id:       id,
		Flange:   dto.Flange,
		TypeFlId: dto.TypeFlId,
		Dn:       dto.Dn,
		Pn:       dto.Pn,
		TypePr:   dto.TypePr,
		StandId:  dto.StandId,
		D4:       dto.D4,
		D3:       dto.D3,
		D2:       dto.D2,
		D1:       dto.D1,
		H:        dto.H,
		S2:       dto.S2,
		S3:       dto.S3,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: size.Id, Message: "Updated"})
}

// @Summary Delete Size
// @Tags Sealur Pro -> sizes
// @Security ApiKeyAuth
// @Description удаление размеров
// @ModuleID deleteSize
// @Accept json
// @Produce json
// @Param id path string true "size id"
// @Param flange query string true "flange"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes/{id} [delete]
func (h *Handler) deleteSize(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	flange := c.Query("flange")
	if flange == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty flange", "empty flange param")
		return
	}

	size, err := h.proClient.DeleteSize(c, &proto.DeleteSizeRequest{
		Id:     id,
		Flange: flange,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: size.Id, Message: "Deleted"})
}

// @Summary Delete Size
// @Tags Sealur Pro -> sizes
// @Security ApiKeyAuth
// @Description удаление всех размеров
// @ModuleID deleteSize
// @Accept json
// @Produce json
// @Param flange query string true "flange"
// @Param typePr query string true "type pr"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/sizes/all [delete]
func (h *Handler) deleteAllSize(c *gin.Context) {
	flange := c.Query("flange")
	if flange == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty flange", "empty flange param")
		return
	}

	typePr := c.Query("typePr")
	if flange == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty type", "empty type param")
		return
	}

	size, err := h.proClient.DeleteAllSize(c, &proto.DeleteAllSizeRequest{
		Flange: flange,
		TypePr: typePr,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: size.Id, Message: "Deleted"})
}
