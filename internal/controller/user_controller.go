package controller

import (
	"net/http"

	"github.com/Meystergod/gochat/internal/apperror"
	"github.com/Meystergod/gochat/internal/domain"
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
		return apperror.NewAppError(apperror.ErrorValidatePayload, err.Error())
	}

	createdUserID, err := userController.userUsecase.Signup(c.Request().Context(), payload.ToModel())
	if err != nil {
		return err
	}

	return utils.Negotiate(c, http.StatusCreated, map[string]string{"id": createdUserID})
}

func (userController *UserController) GetUserInfo(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return apperror.NewAppError(apperror.ErrorGetUrlParams, "could not get user id")
	}

	user, err := userController.userUsecase.GetUserInfo(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return utils.Negotiate(c, http.StatusOK, map[string]domain.User{"user": *user})
}

func (userController *UserController) GetAllUsersInfo(c echo.Context) error {
	users, err := userController.userUsecase.GetAllUsersInfo(c.Request().Context())
	if err != nil {
		return err
	}

	if len(*users) == 0 {
		return utils.Negotiate(c, http.StatusOK, map[string]string{"users": "list is empty"})
	}

	return utils.Negotiate(c, http.StatusOK, map[string][]domain.User{"users": *users})
}

func (userController *UserController) UpdateUserInfo(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	var payload UpdateUserDTO

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return apperror.NewAppError(apperror.ErrorValidatePayload, err.Error())
	}

	user := payload.ToModel()
	user.ID = id

	err := userController.userUsecase.UpdateUserInfo(c.Request().Context(), user)
	if err != nil {
		return err
	}

	return utils.Negotiate(c, http.StatusCreated, map[string]string{"id": id})
}

func (userController *UserController) DeleteUserAccount(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return apperror.NewAppError(apperror.ErrorGetUrlParams, "could not get user id")
	}

	err := userController.userUsecase.DeleteUserAccount(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return utils.Negotiate(c, http.StatusCreated, map[string]string{"id": id})
}
