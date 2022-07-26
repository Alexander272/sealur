package moment

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/moment_model"
	"github.com/Alexander272/sealur_proto/api/moment_api"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initEnvRoutes(api *gin.RouterGroup) {
	env := api.Group("/env", h.middleware.UserIdentity)
	{
		env.GET("/", h.getEnv)
		env = env.Group("/", h.middleware.AccessForMomentAdmin)
		{
			env.POST("/", h.createEnv)
			env.PUT("/:id", h.updateEnv)
			env.DELETE("/:id", h.deleteEnv)
		}
	}
}

// @Summary Get Env
// @Tags Sealur Moment -> env
// @Security ApiKeyAuth
// @Description получение типов сред
// @ModuleID getEnv
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{Data=[]moment_api.Env}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/env/ [get]
func (h *Handler) getEnv(c *gin.Context) {
	env, err := h.gasketClient.GetEnv(c, &moment_api.GetEnvRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: env.Env, Count: len(env.Env)})
}

// @Summary Create Env
// @Tags Sealur Moment -> env
// @Security ApiKeyAuth
// @Description создание среды
// @ModuleID createEnv
// @Accept json
// @Produce json
// @Param env body moment_model.EnvDTO true "env info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/env/ [post]
func (h *Handler) createEnv(c *gin.Context) {
	var dto moment_model.EnvDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	env, err := h.gasketClient.CreateEnv(c, &moment_api.CreateEnvRequest{Title: dto.Title})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-moment/env/%s", env.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: env.Id, Message: "Created"})
}

// @Summary Update Env
// @Tags Sealur Moment -> env
// @Security ApiKeyAuth
// @Description обновление среды
// @ModuleID updateEnv
// @Accept json
// @Produce json
// @Param id path string true "env id"
// @Param env body moment_model.EnvDTO true "env info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/env/{id} [put]
func (h *Handler) updateEnv(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto moment_model.EnvDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.gasketClient.UpdateEnv(c, &moment_api.UpdateEnvRequest{Id: id, Title: dto.Title})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Updated"})
}

// @Summary Delete Env
// @Tags Sealur Moment -> env
// @Security ApiKeyAuth
// @Description Удаление среды
// @ModuleID deleteEnv
// @Accept json
// @Produce json
// @Param id path string true "env id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-moment/env/{id} [delete]
func (h *Handler) deleteEnv(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.gasketClient.DeleteEnv(c, &moment_api.DeleteEnvRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: id, Message: "Deleted"})
}
