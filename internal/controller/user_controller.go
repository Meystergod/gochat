package controller

import (
	"net/http"

	"github.com/Meystergod/gochat/internal/usecase/usecase_user"
	"github.com/Meystergod/gochat/internal/utils"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userUsecase *usecase_user.UserUsecase
}

func NewUserController(userUsecase *usecase_user.UserUsecase) *UserController {
	return &UserController{userUsecase: userUsecase}
}

func (userController *UserController) Signup(c echo.Context) error {
	var payload CreateUserDTO

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorBindAndValidatePayload.Error())
	}

	createdUserID, err := userController.userUsecase.Signup(c.Request().Context(), payload.ToModel())
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusCreated, createdUserID)
}

func (userController *UserController) GetUserInfo(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	user, err := userController.userUsecase.GetUserInfo(c.Request().Context(), id)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, user)
}

func (userController *UserController) GetAllUsersInfo(c echo.Context) error {
	users, err := userController.userUsecase.GetAllUsersInfo(c.Request().Context())
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	if len(*users) == 0 {
		return utils.Negotiate(c, http.StatusOK, "list is empty")
	}

	return utils.Negotiate(c, http.StatusOK, users)
}

func (userController *UserController) UpdateUserInfo(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	var payload UpdateUserDTO

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorBindAndValidatePayload.Error())
	}

	user := payload.ToModel()
	user.ID = id

	err := userController.userUsecase.UpdateUserInfo(c.Request().Context(), user)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, id)
}

func (userController *UserController) DeleteUserAccount(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	err := userController.userUsecase.DeleteUserAccount(c.Request().Context(), id)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusNoContent, id)
}
