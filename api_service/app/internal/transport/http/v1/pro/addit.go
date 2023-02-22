package pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/pro_model"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initAdditRoutes(api *gin.RouterGroup) {
	addit := api.Group("/additionals")
	{
		addit.GET("/", h.getAddit)
		addit = addit.Group("", h.middleware.UserIdentity, h.middleware.AccessForProAdmin)
		{
			addit.POST("/", h.createAddit)
			addit.PATCH("/:id/mat", h.updateMat)
			addit.PATCH("/:id/mod", h.updateMod)
			addit.PATCH("/:id/temp", h.updateTemp)
			addit.PATCH("/:id/moun", h.updateMoun)
			addit.PATCH("/:id/grap", h.updateGrap)
			addit.PATCH("/:id/fil", h.updateFillers)
			addit.PATCH("/:id/coat", h.updateCoating)
			addit.PATCH("/:id/constr", h.updateConstruction)
			addit.PATCH("/:id/obt", h.updateObturator)
			addit.PATCH("/:id/basis", h.updateBasis)
			addit.PATCH("/:id/pobt", h.updatePObturator)
			addit.PATCH("/:id/seal", h.updateSealant)
		}
	}
}

// @Summary Get Addit
// @Tags Sealur Pro -> additionals
// @Description получение всех дополнительных сведенний (список материалов, креплений и тд)
// @ModuleID getAddit
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{data=[]pro_api.Additional}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals [get]
func (h *Handler) getAddit(c *gin.Context) {
	add, err := h.proClient.GetAllAdditional(c, &pro_api.GetAllAddRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: add.Additionals})
}

// @Summary Create Addit
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description создание дополнительных сведенний
// @ModuleID createAddit
// @Accept json
// @Produce json
// @Param data body pro_model.CreateAdditDTO true "addit info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals [post]
func (h *Handler) createAddit(c *gin.Context) {
	var dto pro_model.CreateAdditDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.proClient.CreateAdditional(c, &pro_api.CreateAddRequest{
		Materials:    dto.Materials,
		Mod:          dto.Mod,
		Temperature:  dto.Temperature,
		Mounting:     dto.Mounting,
		Graphite:     dto.Graphite,
		Fillers:      dto.Fillers,
		Coating:      dto.Coating,
		Construction: dto.Construction,
		Obturator:    dto.Obturator,
		Basis:        dto.Basis,
		Sealant:      dto.Sealant,
		PObturator:   dto.PObt,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

// @Summary Update Mat
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление материалов
// @ModuleID updateMat
// @Accept json
// @Produce json
// @Param data body pro_model.UpdateMatDTO true "additional materials info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/mat [patch]
func (h *Handler) updateMat(c *gin.Context) {
	var dto pro_model.UpdateMatDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateMat(c, &pro_api.UpdateAddMatRequest{Id: id, Materials: dto.Materials, TypeCh: dto.TypeCh, Change: dto.Change})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated materials"})
}

// @Summary Update Mod
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление модифицирующего элемента
// @ModuleID updateMod
// @Accept json
// @Produce json
// @Param data body pro_model.UpdateModDTO true "additional modification info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/mod [patch]
func (h *Handler) updateMod(c *gin.Context) {
	var dto pro_model.UpdateModDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateMod(c, &pro_api.UpdateAddModRequest{Id: id, Mod: dto.Mod, TypeCh: dto.TypeCh, Change: dto.Change})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated modifications"})
}

// @Summary Update Temp
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление температур
// @ModuleID updateTemp
// @Accept json
// @Produce json
// @Param data body pro_model.UpdateTempDTO true "additional temperature info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/temp [patch]
func (h *Handler) updateTemp(c *gin.Context) {
	var dto pro_model.UpdateTempDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateTemp(c, &pro_api.UpdateAddTemRequest{Id: id, Temperature: dto.Temperature, TypeCh: dto.TypeCh, Change: dto.Change})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated temperature"})
}

// @Summary Update Moun
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление крепления на вертикальном фланце
// @ModuleID updateMoun
// @Accept json
// @Produce json
// @Param data body pro_model.UpdateMounDTO true "additional mounting info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/moun [patch]
func (h *Handler) updateMoun(c *gin.Context) {
	var dto pro_model.UpdateMounDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateMoun(c, &pro_api.UpdateAddMounRequest{Id: id, Mounting: dto.Mounting, TypeCh: dto.TypeCh, Change: dto.Change})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated mounting"})
}

// @Summary Update Grap
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление графита
// @ModuleID updateGrap
// @Accept json
// @Produce json
// @Param data body pro_model.UpdateGrapDTO true "additional graphite info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/grap [patch]
func (h *Handler) updateGrap(c *gin.Context) {
	var dto pro_model.UpdateGrapDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateGrap(c, &pro_api.UpdateAddGrapRequest{Id: id, Graphite: dto.Graphite, TypeCh: dto.TypeCh, Change: dto.Change})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated graphite"})
}

// @Summary Update Fillers
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление наполнителя для снп
// @ModuleID updateFillers
// @Accept json
// @Produce json
// @Param data body pro_model.UpdateFillersDTO true "additional fillers info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/fil [patch]
func (h *Handler) updateFillers(c *gin.Context) {
	var dto pro_model.UpdateFillersDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateFillers(c, &pro_api.UpdateAddFillersRequest{Id: id, Fillers: dto.Fillers, TypeCh: dto.TypeCh, Change: dto.Change})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated fillers"})
}

// @Summary Update Coating
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление способа исполнения для путг и путгм
// @ModuleID updateCoating
// @Accept json
// @Produce json
// @Param data body pro_model.UpdateCoatingDTO true "additional coating info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/coat [patch]
func (h *Handler) updateCoating(c *gin.Context) {
	var dto pro_model.UpdateCoatingDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateCoating(c, &pro_api.UpdateAddCoatingRequest{Id: id, Coating: dto.Coating, TypeCh: dto.TypeCh, Change: dto.Change})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated coating"})
}

// @Summary Update Construction
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление конструкций для путг
// @ModuleID updateConstruction
// @Accept json
// @Produce json
// @Param data body pro_model.UpdateConstrDTO true "additional construction info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/constr [patch]
func (h *Handler) updateConstruction(c *gin.Context) {
	var dto pro_model.UpdateConstrDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateConstruction(c, &pro_api.UpdateAddConstructionRequest{Id: id, Constr: dto.Constr, TypeCh: dto.TypeCh, Change: dto.Change})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated construction"})
}

// @Summary Update Obturator
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление обтюраторов для путг
// @ModuleID updateObturator
// @Accept json
// @Produce json
// @Param data body pro_model.UpdateObturatorDTO true "additional obturation info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/obt [patch]
func (h *Handler) updateObturator(c *gin.Context) {
	var dto pro_model.UpdateObturatorDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateObturator(c, &pro_api.UpdateAddObturatorRequest{Id: id, Obturator: dto.Obt, TypeCh: dto.TypeCh, Change: dto.Change})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated obturators"})
}

// @Summary Update Basis
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление конструкций для путгм
// @ModuleID updateBasis
// @Accept json
// @Produce json
// @Param data body pro_model.UpdateBasisDTO true "additional basis info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/basis [patch]
func (h *Handler) updateBasis(c *gin.Context) {
	var dto pro_model.UpdateBasisDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateBasis(c, &pro_api.UpdateAddBasisRequest{Id: id, Basis: dto.Basis, TypeCh: dto.TypeCh, Change: dto.Change})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated basis"})
}

// @Summary Update PObturator
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление обтюраторов для путгм
// @ModuleID updatePObturator
// @Accept json
// @Produce json
// @Param data body pro_model.UpdatePObtDTO true "additional p_obturator info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/pobt [patch]
func (h *Handler) updatePObturator(c *gin.Context) {
	var dto pro_model.UpdatePObtDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdatePObturator(c, &pro_api.UpdateAddPObturatorRequest{Id: id, PObturator: dto.Obturator, TypeCh: dto.TypeCh, Change: dto.Change})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated sealant"})
}

// @Summary Update Sealant
// @Tags Sealur Pro -> additionals
// @Security ApiKeyAuth
// @Description обновление уплотнителя для путгм
// @ModuleID updateSealant
// @Accept json
// @Produce json
// @Param data body pro_model.UpdateSealantDTO true "additional sealant info"
// @Param id path string true "addit id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/additionals/{id}/seal [patch]
func (h *Handler) updateSealant(c *gin.Context) {
	var dto pro_model.UpdateSealantDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.proClient.UpdateSealant(c, &pro_api.UpdateAddSealantRequest{Id: id, Sealant: dto.Sealant, TypeCh: dto.TypeCh, Change: dto.Change})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated sealant"})
}
