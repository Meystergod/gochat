package dto

import "github.com/google/uuid"

type CreateUserDTO struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type UpdateUserDTO struct {
	ID       uuid.UUID `json:"id" bson:"_id,omitempty"`
	Name     string    `json:"name" bson:"name"`
	Email    string    `json:"email" bson:"email"`
	Password string    `json:"password" bson:"password"`
}
