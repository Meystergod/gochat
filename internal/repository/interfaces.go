package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Meystergod/gochat/internal/entity/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (primitive.ObjectID, error)
	GetUser(ctx context.Context, oid primitive.ObjectID) (*model.User, error)
	GetAllUsers(ctx context.Context) (*[]model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, oid primitive.ObjectID) error
}
