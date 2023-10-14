package repository

import (
	"context"

	"github.com/Meystergod/gochat/internal/entity/model"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (string, error)
	GetUser(ctx context.Context, uuid uuid.UUID) (*model.User, error)
	GetAllUsers(ctx context.Context) (*[]model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, uuid uuid.UUID) error
}
