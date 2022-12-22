package user

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/user_model"
	"github.com/Alexander272/sealur_proto/api/user_api"
	"github.com/gin-gonic/gin"
)

// @Summary Create Role
// @Tags Sealur Pro -> users
// @Security ApiKeyAuth
// @Description добавление роли для пользователя
// @ModuleID createRole
// @Accept json
// @Produce json
// @Param role body user_model.UserRole true "role info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /users/roles [post]
func (h *Handler) createRole(c *gin.Context) {
	var dto user_model.UserRole
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}
	if dto.UserId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "user id is empty", "invalid data send")
		return
	}

	req := user_api.CreateRoleRequest{
		UserId:  dto.UserId,
		Service: dto.Service,
		Role:    dto.Role,
	}

	_, err := h.userClient.CreateRole(c, &req)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "User role added successfully"})
}

// @Summary Update Role
// @Tags Sealur Pro -> users
// @Security ApiKeyAuth
// @Description обновление роли пользователя
// @ModuleID updateRole
// @Accept json
// @Produce json
// @Param id path string true "role id"
// @Param role body user_model.UserRole true "role info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /users/roles/{id} [put]
func (h *Handler) updateRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto user_model.UserRole
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	if err := h.services.Session.SingOut(c, dto.UserId); err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "failed to close session")
		return
	}

	req := user_api.UpdateRoleRequest{
		Id:      id,
		Service: dto.Service,
		Role:    dto.Role,
	}

	role, err := h.userClient.UpdateRole(c, &req)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: role.Id, Message: "User role updated successfully"})
}

// @Summary Delete Role
// @Tags Sealur Pro -> users
// @Security ApiKeyAuth
// @Description удаление роли пользователя
// @ModuleID deleteRole
// @Accept json
// @Produce json
// @Param id path string true "role id"
// @Param userId query string true "user id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /users/roles/{id} [delete]
func (h *Handler) deleteRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	userId := c.Query("userId")
	if userId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty userId", "empty userId param")
		return
	}

	if err := h.services.Session.SingOut(c, userId); err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "failed to close session")
		return
	}

	role, err := h.userClient.DeleteRole(c, &user_api.DeleteRoleRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: role.Id, Message: "User role deleted successfully"})
}
