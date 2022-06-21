package user

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto/proto_user"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/gin-gonic/gin"
)

// @Summary SignIn
// @Tags auth
// @Description вход в систему
// @ModuleID signIn
// @Accept json
// @Produce json
// @Param data body models.SignIn true "credentials"
// @Success 200 {object} structs.Tokens
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var dto models.SignIn
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.userClient.GetUser(c, &proto_user.GetUserRequest{Login: dto.Login, Password: dto.Password})
	if err != nil {
		logger.Debug("err: ", err.Error())
		// if errors.Is(err, models.ErrFlangeAlreadyExists) {
		// 	models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
		// 	return
		// }
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	//TODO добавить куки и запись в редисе

	c.JSON(http.StatusOK, models.IdResponse{Message: "Authorization completed successfully"})
}

// @Summary SignUp
// @Tags auth
// @Description регистрация
// @ModuleID singUp
// @Accept json
// @Produce json
// @Param data body models.SignUp true "user info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /auth/sign-up [post]
func (h *Handler) singUp(c *gin.Context) {
	var dto models.SignUp
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	req := proto_user.CreateUserRequest{
		Organization: dto.Organization,
		Name:         dto.Name,
		Email:        dto.Email,
		City:         dto.City,
		Position:     dto.Position,
		Phone:        dto.Phone,
	}

	_, err := h.userClient.CreateUser(c, &req)
	if err != nil {
		// if errors.Is(err, models.ErrFlangeAlreadyExists) {
		// 	models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
		// 	return
		// }
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Registration completed successfully"})
}

// @Summary SignOut
// @Tags auth
// @Description выход из аккаунта
// @ModuleID signOut
// @Accept json
// @Produce json
// @Success 204 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /auth/sign-out [post]
func (h *Handler) signOut(c *gin.Context) {
	//TODO удалить куки и запись в редисе

	c.JSON(http.StatusNoContent, models.IdResponse{Message: "Sign-out completed successfully"})
}
