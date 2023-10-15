package usecase

import (
	"context"
	"time"

	"github.com/Meystergod/gochat/internal/entity/model"
	"github.com/Meystergod/gochat/internal/repository"
	"github.com/Meystergod/gochat/internal/utils"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUsecase struct {
	userRepository repository.UserRepository
}

func NewUserUsecase(userRepository repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
	}
}

func (userUsecase *UserUsecase) CreateUser(ctx context.Context, user *model.User) (string, error) {
	user.LastVisitAt = time.Now()
	user.RegisteredAt = time.Now()

	oid, err := userUsecase.userRepository.CreateUser(ctx, user)
	if err != nil {
		return utils.EmptyString, err
	}

	return oid.Hex(), nil
}

func (userUsecase *UserUsecase) GetUserInfo(ctx context.Context, id string) (*model.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, utils.ErrorConvert.Error())
	}

	user, err := userUsecase.userRepository.GetUser(ctx, oid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userUsecase *UserUsecase) GetAllUsersInfo(ctx context.Context) (*[]model.User, error) {
	users, err := userUsecase.userRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (userUsecase *UserUsecase) UpdateUserInfo(ctx context.Context, user *model.User, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, utils.ErrorConvert.Error())
	}

	user.ID = oid

	err = userUsecase.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (userUsecase *UserUsecase) DeleteUserAccount(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(errors.New("error convert hex to oid"), utils.ErrorConvert.Error())
	}

	err = userUsecase.userRepository.DeleteUser(ctx, oid)
	if err != nil {
		return err
	}

	return nil
}
