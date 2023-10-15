package dto

import (
	"github.com/Meystergod/gochat/internal/entity/model"
)

type CreateUserDTO struct {
	Name     string `json:"name" bson:"name" validate:"required"`
	Email    string `json:"email" bson:"email" validate:"required"`
	Password string `json:"password" bson:"password" validate:"required"`
}

type UpdateUserDTO struct {
	Name     string `json:"name" bson:"name" validate:"required"`
	Email    string `json:"email" bson:"email" validate:"required"`
	Password string `json:"password" bson:"password" validate:"required"`
}

func (createUserDTO *CreateUserDTO) ToModel() *model.User {
	return &model.User{
		Name:     createUserDTO.Name,
		Email:    createUserDTO.Email,
		Password: createUserDTO.Password,
	}
}

func (updateUserDTO *UpdateUserDTO) ToModel() *model.User {
	return &model.User{
		Name:     updateUserDTO.Name,
		Email:    updateUserDTO.Email,
		Password: updateUserDTO.Password,
	}
}
