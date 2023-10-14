package usecase

import (
	"github.com/Meystergod/gochat/internal/controller"
)

type UserUsecase struct {
	userController *controller.UserController
}

func NewUserUsecase(userController *controller.UserController) *UserUsecase {
	return &UserUsecase{
		userController: userController,
	}
}

func (userUsecase *UserUsecase) Login() error {
	return nil
}
