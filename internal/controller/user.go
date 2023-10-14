package controller

import (
	"net/http"

	"github.com/Meystergod/gochat/internal/entity/dto"
	"github.com/Meystergod/gochat/internal/repository"
	"github.com/Meystergod/gochat/internal/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	userRepository repository.UserRepository
}

func NewUserController(userRepository repository.UserRepository) *UserController {
	return &UserController{userRepository: userRepository}
}

func (userController *UserController) CreateUser(c echo.Context) error {
	var payload dto.CreateUserDTO

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorBindAndValidatePayload.Error())
	}

	_, err := userController.userRepository.CreateUser(c.Request().Context(), payload.ToModel())
	if err == nil {
		return utils.Negotiate(c, http.StatusConflict, "user with this title is exist")
	}

	createdUserID, err := userController.userRepository.CreateUser(c.Request().Context(), payload.ToModel())
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusCreated, createdUserID)
}

func (userController *UserController) GetUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	userUUID, err := uuid.FromBytes([]byte(id))
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	user, err := userController.userRepository.GetUser(c.Request().Context(), userUUID)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, user)
}

func (userController *UserController) GetAllUsers(c echo.Context) error {
	users, err := userController.userRepository.GetAllUsers(c.Request().Context())
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, users)
}

func (userController *UserController) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	var payload dto.UpdateUserDTO

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorBindAndValidatePayload.Error())
	}

	user := payload.ToModel()

	userUUID, err := uuid.FromBytes([]byte(id))
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	user.ID = userUUID

	err = userController.userRepository.UpdateUser(c.Request().Context(), user)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, nil)
}

func (userController *UserController) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	userUUID, err := uuid.FromBytes([]byte(id))
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	err = userController.userRepository.DeleteUser(c.Request().Context(), userUUID)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusNoContent, nil)
}
