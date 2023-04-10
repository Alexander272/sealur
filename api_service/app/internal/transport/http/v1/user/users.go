package user

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/email_api"
	"github.com/Alexander272/sealur_proto/api/user/user_api"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userApi    user_api.UserServiceClient
	emailApi   email_api.EmailServiceClient
	auth       config.AuthConfig
	http       config.HttpConfig
	services   *service.Services
	cookieName string
}

func NewUserHandler(
	userApi user_api.UserServiceClient, emailApi email_api.EmailServiceClient,
	auth config.AuthConfig, http config.HttpConfig,
	services *service.Services,
	cookieName string,
) *UserHandler {
	return &UserHandler{
		userApi:    userApi,
		emailApi:   emailApi,
		auth:       auth,
		services:   services,
		cookieName: cookieName,
	}
}

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	handler := NewUserHandler(h.userApi, h.emailApi, h.auth, h.http, h.services, h.cookieName)

	users := api.Group("/users")
	{
		users.GET("/managers", handler.getManagers)
		users.GET("/:id", handler.getUser)
		users.POST("/confirm/:code", handler.confirm)
		users.POST("/manager", handler.setManager)
		users.POST("/recovery", handler.recoveryPassword)
		users.POST("/recovery/:code", handler.recoveryPasswordCode)
	}
}

func (h *UserHandler) getManagers(c *gin.Context) {
	users, err := h.userApi.GetManagers(c, &user_api.GetNewUser{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: users})
}

func (h *UserHandler) getUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	user, err := h.userApi.Get(c, &user_api.GetUser{Id: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: user})
}

func (h *UserHandler) confirm(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty code", "empty code param")
		return
	}

	data, err := h.services.Confirm.Get(c, code)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	user, err := h.userApi.Confirm(c, &user_api.ConfirmUser{Id: data.UserId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	// запись в редисе и генерация токенов
	token, err := h.services.SignIn(c, user)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.SetCookie(h.cookieName, token, int(h.auth.RefreshTokenTTL.Seconds()), "/", c.Request.Host, h.auth.Secure, true)
	c.JSON(http.StatusOK, models.DataResponse{Data: user})
}

func (h *UserHandler) setManager(c *gin.Context) {
	var dto *user_api.UserManager
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.userApi.SetManager(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated manager successfully"})
}

func (h *UserHandler) recoveryPassword(c *gin.Context) {
	var dto *user_api.GetUserByEmail
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Введены некорректные данные")
		return
	}

	user, err := h.userApi.GetByEmail(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка")
		return
	}

	// генерировать код для подтверждения и записывать его в редис (с id пользователя)
	code, err := h.services.Confirm.Create(c, user.Id)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка")
		return
	}

	data := &email_api.RecoveryPassword{
		Email: user.Email,
		//TODO использовать тут хост как-то не особо правильно выглядит
		// Link: fmt.Sprintf("%s/auth/recovery/%s", h.http.Host, code),
		Link: fmt.Sprintf("%s/auth/recovery/%s", "http://pro.sealur.ru", code),
	}

	_, err = h.emailApi.Recovery(c, data)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка при отправлении письма")
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Reset completed successfully"})
}

func (h *UserHandler) recoveryPasswordCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty code", "Некорректный код восстановления")
		return
	}

	var dto *user_api.GetUserByEmail
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Введены некорректные данные")
		return
	}

	data, err := h.services.Confirm.Get(c, code)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка")
		return
	}

	user, err := h.userApi.Get(c, &user_api.GetUser{Id: data.UserId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка")
		return
	}

	_, err = h.userApi.Update(c, &user_api.UpdateUser{Id: user.Id, Password: dto.Password})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка")
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated password successfully"})
}

// // @Summary Get All Users
// // @Tags Users
// // @Security ApiKeyAuth
// // @Description получение всех пользователей
// // @ModuleID getAllUser
// // @Accept json
// // @Produce json
// // @Param page path string true "page number"
// // @Param search query string false "search value"
// // @Param search_field query false "search field"
// // @Success 200 {object} models.DataResponse{data=[]user_api.User}
// // @Failure 400,404 {object} models.ErrorResponse
// // @Failure 500 {object} models.ErrorResponse
// // @Failure default {object} models.ErrorResponse
// // @Router /users/all/{page} [get]
// func (h *Handler) getAllUsers(c *gin.Context) {
// 	var countPage, limit int32
// 	page := c.Param("page")
// 	if page == "" {
// 		countPage = 1
// 	} else {
// 		count, err := strconv.Atoi(page)
// 		if err != nil {
// 			models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "something went wrong")
// 			return
// 		}
// 		countPage = int32(count)
// 	}

// 	parLimit := c.Query("limit")
// 	if parLimit == "" {
// 		limit = 12
// 	} else {
// 		l, err := strconv.Atoi(parLimit)
// 		if err != nil {
// 			logger.Debug("err: ", err.Error())
// 			models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "something went wrong")
// 			return
// 		}
// 		limit = int32(l)
// 	}

// 	search := c.Query("search")
// 	Fields := []user_api.GetAllUserRequest_SearchField{}
// 	if search != "" {
// 		fields := c.QueryArray("search_field")
// 		for _, v := range fields {
// 			field := user_api.GetAllUserRequest_SearchField_value[v]
// 			Fields = append(Fields, user_api.GetAllUserRequest_SearchField(field))
// 		}
// 	}

// 	users, err := h.userClient.GetAllUsers(c, &user_api.GetAllUserRequest{
// 		Limit:  limit,
// 		Offset: limit * (countPage - 1),
// 		Search: &user_api.GetAllUserRequest_Search{
// 			Field: Fields,
// 			Value: search,
// 		},
// 	})
// 	if err != nil {
// 		logger.Debug("err: ", err.Error())
// 		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.DataResponse{Data: users.Users, Count: int(math.Ceil(float64(users.Count) / float64(limit)))})
// }

// // @Summary Get New Users
// // @Tags Users
// // @Security ApiKeyAuth
// // @Description получение новых пользователей
// // @ModuleID getNewUsers
// // @Accept json
// // @Produce json
// // @Success 200 {object} models.DataResponse{data=[]user_api.User}
// // @Failure 400,404 {object} models.ErrorResponse
// // @Failure 500 {object} models.ErrorResponse
// // @Failure default {object} models.ErrorResponse
// // @Router /users/new [get]
// func (h *Handler) getNewUsers(c *gin.Context) {
// 	users, err := h.userClient.GetNewUsers(c, &user_api.GetNewUserRequest{})
// 	if err != nil {
// 		logger.Debug("err: ", err.Error())
// 		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.DataResponse{Data: users.Users, Count: len(users.Users)})
// }

// // @Summary Get User
// // @Tags Users
// // @Security ApiKeyAuth
// // @Description получение данных пользователя
// // @ModuleID getUser
// // @Accept json
// // @Produce json
// // @Param id path string true "user id"
// // @Success 200 {object} models.DataResponse{data=user_api.User}
// // @Failure 400,404 {object} models.ErrorResponse
// // @Failure 500 {object} models.ErrorResponse
// // @Failure default {object} models.ErrorResponse
// // @Router /users/{id} [get]
// func (h *Handler) getUser(c *gin.Context) {
// 	id := c.Param("id")
// 	if id == "" {
// 		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
// 		return
// 	}

// 	user, err := h.userClient.GetUser(c, &user_api.GetUserRequest{Id: id})
// 	if err != nil {
// 		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.DataResponse{Data: user.User})
// }

// // @Summary Confirm User
// // @Tags Users
// // @Security ApiKeyAuth
// // @Description подтверждение пользователя
// // @ModuleID confirmUser
// // @Accept json
// // @Produce json
// // @Param user body user_model.ConfirmUser true "user info"
// // @Success 200 {object} models.IdResponse
// // @Failure 400,404 {object} models.ErrorResponse
// // @Failure 500 {object} models.ErrorResponse
// // @Failure default {object} models.ErrorResponse
// // @Router /users/confirm [post]
// func (h *Handler) confirmUser(c *gin.Context) {
// 	var dto user_model.ConfirmUser
// 	if err := c.BindJSON(&dto); err != nil {
// 		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
// 		return
// 	}

// 	var roles []*user_api.Role
// 	for _, ur := range dto.Roles {
// 		roles = append(roles, &user_api.Role{
// 			Service: ur.Service,
// 			Role:    ur.Role,
// 		})
// 	}

// 	req := user_api.ConfirmUserRequest{
// 		Id:       dto.Id,
// 		Login:    dto.Login,
// 		Password: dto.Password,
// 		Roles:    roles,
// 	}

// 	_, err := h.userClient.ConfirmUser(c, &req)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "user already exists") {
// 			models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "User with this login already exists")
// 			return
// 		}

// 		if strings.Contains(err.Error(), "failed to send") {
// 			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Failed to send email")
// 			return
// 		}

// 		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.IdResponse{Message: "User successfully verified"})
// }

// // @Summary Reject User
// // @Tags Users
// // @Security ApiKeyAuth
// // @Description отклонение пользователя
// // @ModuleID rejectUser
// // @Accept json
// // @Produce json
// // @Success 200 {object} models.IdResponse
// // @Failure 400,404 {object} models.ErrorResponse
// // @Failure 500 {object} models.ErrorResponse
// // @Failure default {object} models.ErrorResponse
// // @Router /users/reject/{id} [delete]
// func (h *Handler) rejectUser(c *gin.Context) {
// 	id := c.Param("id")
// 	if id == "" {
// 		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
// 		return
// 	}

// 	req := user_api.DeleteUserRequest{
// 		Id: id,
// 	}

// 	_, err := h.userClient.RejectUser(c, &req)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "failed to send") {
// 			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Failed to send email")
// 			return
// 		}

// 		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.IdResponse{Message: "User rejected"})
// }

// // @Summary Update User
// // @Tags Users
// // @Security ApiKeyAuth
// // @Description обновление данных пользователя
// // @ModuleID updateUser
// // @Accept json
// // @Produce json
// // @Param id path string true "user id"
// // @Param user body user_model.UpdateUser true "user info"
// // @Success 200 {object} models.IdResponse
// // @Failure 400,404 {object} models.ErrorResponse
// // @Failure 500 {object} models.ErrorResponse
// // @Failure default {object} models.ErrorResponse
// // @Router /users/{id} [patch]
// func (h *Handler) updateUser(c *gin.Context) {
// 	id := c.Param("id")
// 	if id == "" {
// 		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
// 		return
// 	}
// 	var dto user_model.UpdateUser
// 	if err := c.BindJSON(&dto); err != nil {
// 		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
// 		return
// 	}

// 	req := user_api.UpdateUserRequest{
// 		Id:       id,
// 		Name:     dto.Name,
// 		Email:    dto.Email,
// 		Position: dto.Position,
// 		Phone:    dto.Phone,
// 		Login:    dto.Login,
// 		Password: dto.Password,
// 	}

// 	_, err := h.userClient.UpdateUser(c, &req)
// 	if err != nil {
// 		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.IdResponse{Message: "User data updated successfully"})
// }

// // @Summary Delete User
// // @Tags Users
// // @Security ApiKeyAuth
// // @Description удаление пользователя
// // @ModuleID deleteUser
// // @Accept json
// // @Produce json
// // @Param id path string true "user id"
// // @Success 200 {object} models.IdResponse
// // @Failure 400,404 {object} models.ErrorResponse
// // @Failure 500 {object} models.ErrorResponse
// // @Failure default {object} models.ErrorResponse
// // @Router /users/{id} [delete]
// func (h *Handler) deleteUser(c *gin.Context) {
// 	id := c.Param("id")
// 	if id == "" {
// 		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
// 		return
// 	}

// 	if err := h.services.Session.SingOut(c, id); err != nil {
// 		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "failed to close session")
// 		return
// 	}

// 	_, err := h.userClient.DeleteUser(c, &user_api.DeleteUserRequest{Id: id})
// 	if err != nil {
// 		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.IdResponse{Message: "User deleted successfully"})
// }

// // @Summary Clear Limit
// // @Tags Auth
// // @Description сброс счетчика неуспешных попыток авторизации
// // @ModuleID clearLimit
// // @Accept json
// // @Produce json
// // @Success 200 {object} models.StatusResponse
// // @Failure 400,404 {object} models.ErrorResponse
// // @Failure 500 {object} models.ErrorResponse
// // @Failure default {object} models.ErrorResponse
// // @Router /users/clear [post]
// func (h *Handler) clearLimit(c *gin.Context) {
// 	ip := c.Query("ip")
// 	if ip == "" {
// 		models.NewErrorResponse(c, http.StatusBadRequest, "ip is empty", "ip is empty")
// 		return
// 	}

// 	h.services.Remove(c, ip)

// 	c.JSON(http.StatusOK, models.StatusResponse{Status: "Success"})
// }
