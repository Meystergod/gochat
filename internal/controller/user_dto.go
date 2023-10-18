package controller

import "github.com/Meystergod/gochat/internal/domain"

type CreateUserDTO struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserDTO struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (createUserDTO *CreateUserDTO) ToModel() *domain.User {
	return &domain.User{
		Name:     createUserDTO.Name,
		Email:    createUserDTO.Email,
		Password: createUserDTO.Password,
	}
}

func (updateUserDTO *UpdateUserDTO) ToModel() *domain.User {
	return &domain.User{
		Name:     updateUserDTO.Name,
		Email:    updateUserDTO.Email,
		Password: updateUserDTO.Password,
	}
}
