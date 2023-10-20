package controller

import "github.com/Meystergod/gochat/internal/domain"

type CreateUserDTO struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserDTO struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
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
