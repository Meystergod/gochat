package controller

import (
	"net/http"

	"github.com/Meystergod/gochat/internal/entity/dto"
	"github.com/Meystergod/gochat/internal/usecase"
	"github.com/Meystergod/gochat/internal/utils"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userUsecase *usecase.UserUsecase
}

func NewUserController(userUsecase *usecase.UserUsecase) *UserController {
	return &UserController{userUsecase: userUsecase}
}

func (userController *UserController) CreateUser(c echo.Context) error {
	var payload dto.CreateUserDTO

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorBindAndValidatePayload.Error())
	}

	createdUserID, err := userController.userUsecase.CreateUser(c.Request().Context(), payload.ToModel())
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

	var payload dto.UpdateUserDTO

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorBindAndValidatePayload.Error())
	}

	user := payload.ToModel()

	err := userController.userUsecase.UpdateUserInfo(c.Request().Context(), user, id)
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
