package repository_user

import (
	"github.com/Meystergod/gochat/internal/domain"
	"github.com/Meystergod/gochat/internal/utils"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const MethodCreate = "create"
const MethodUpdate = "update"

func userToDomain(u *User) domain.User {
	return domain.User{
		ID:           u.ID.Hex(),
		Name:         u.Name,
		Email:        u.Email,
		Password:     u.Password,
		RegisteredAt: u.RegisteredAt,
	}
}

func userToRepository(user *domain.User, method string) (User, error) {
	switch method {
	case MethodCreate:
		return User{
			Name:         user.Name,
			Email:        user.Email,
			Password:     user.Password,
			RegisteredAt: user.RegisteredAt,
		}, nil
	case MethodUpdate:
		oid, err := primitive.ObjectIDFromHex(user.ID)
		if err != nil {
			return User{}, errors.Wrap(err, utils.ErrorConvert.Error())
		}
		return User{
			ID:       oid,
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
		}, nil
	default:
		return User{}, errors.Wrap(errors.New("unknown method for convert to repository model"), utils.ErrorConvert.Error())
	}
}
