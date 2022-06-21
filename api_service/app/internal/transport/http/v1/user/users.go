package user

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto/proto_user"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/gin-gonic/gin"
)

// @Summary Get All Users
// @Tags Sealur Pro -> users
// @Security ApiKeyAuth
// @Description получение всех пользователей
// @ModuleID getAllUser
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{data=[]proto_user.User}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /users/all [get]
func (h *Handler) getAllUsers(c *gin.Context) {
	users, err := h.userClient.GetAllUsers(c, &proto_user.GetAllUserRequest{})
	if err != nil {
		logger.Debug("err: ", err.Error())
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: users.Users, Count: len(users.Users)})
}

// @Summary Get New Users
// @Tags Sealur Pro -> users
// @Security ApiKeyAuth
// @Description получение новых пользователей
// @ModuleID getNewUsers
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{data=[]proto_user.User}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /users/new [get]
func (h *Handler) getNewUsers(c *gin.Context) {
	users, err := h.userClient.GetNewUsers(c, &proto_user.GetNewUserRequest{})
	if err != nil {
		logger.Debug("err: ", err.Error())
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: users.Users, Count: len(users.Users)})
}

// @Summary Get User
// @Tags Sealur Pro -> users
// @Security ApiKeyAuth
// @Description получение данных пользователя
// @ModuleID getUser
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} models.DataResponse{data=proto_user.User}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /users/{id} [get]
func (h *Handler) getUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	user, err := h.userClient.GetUser(c, &proto_user.GetUserRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: user})
}

// @Summary Confirm User
// @Tags Sealur Pro -> users
// @Security ApiKeyAuth
// @Description потверждение пользователя
// @ModuleID confirmUser
// @Accept json
// @Produce json
// @Param user body models.ConfirmUser true "user info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /users/confirm [post]
func (h *Handler) confirmUser(c *gin.Context) {
	var dto models.ConfirmUser
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	var roles []*proto_user.Role
	for _, ur := range dto.Roles {
		roles = append(roles, &proto_user.Role{
			Service: ur.Service,
			Role:    ur.Role,
		})
	}

	req := proto_user.ConfirmUserRequest{
		Id:       dto.Id,
		Login:    dto.Login,
		Password: dto.Password,
		Roles:    roles,
	}

	_, err := h.userClient.ConfirmUser(c, &req)
	if err != nil {
		//TODO надо отдельно обрабатывать ошибку отправки email
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "User successfully verified"})
}

// @Summary Update User
// @Tags Sealur Pro -> users
// @Security ApiKeyAuth
// @Description обновление данных пользователя
// @ModuleID updateUser
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Param user body models.UpdateUser true "user info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /users/{id} [patch]
func (h *Handler) updateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}
	var dto models.UpdateUser
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	req := proto_user.UpdateUserRequest{
		Id:       id,
		Name:     dto.Name,
		Email:    dto.Email,
		Position: dto.Position,
		Phone:    dto.Phone,
	}

	_, err := h.userClient.UpdateUser(c, &req)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "User data updated successfully"})
}

// @Summary Delete User
// @Tags Sealur Pro -> users
// @Security ApiKeyAuth
// @Description удаление пользователя
// @ModuleID deleteUser
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /users/{id} [delete]
func (h *Handler) deleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	_, err := h.userClient.DeleteUser(c, &proto_user.DeleteUserRequest{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "User deleted successfully"})
}
