package repository_user

import (
	"github.com/Meystergod/gochat/internal/domain"
	"github.com/Meystergod/gochat/internal/utils"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func userToDomain(u *User) domain.User {
	return domain.User{
		ID:           u.ID.Hex(),
		Name:         u.Name,
		Email:        u.Email,
		Password:     u.Password,
		RegisteredAt: u.RegisteredAt,
		LastVisitAt:  u.LastVisitAt,
	}
}

func userToRepository(user *domain.User) (User, error) {
	oid, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return User{}, errors.Wrap(err, utils.ErrorConvert.Error())
	}

	return User{
		ID:           oid,
		Name:         user.Name,
		Email:        user.Email,
		Password:     user.Password,
		RegisteredAt: user.RegisteredAt,
		LastVisitAt:  user.LastVisitAt,
	}, nil
}
