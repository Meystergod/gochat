package usecase_user

import (
	"context"
	"time"

	"github.com/Meystergod/gochat/internal/domain"
	"github.com/Meystergod/gochat/internal/utils"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (string, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	GetAllUsers(ctx context.Context) (*[]domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id string) error
}

type UserUsecase struct {
	userRepository UserRepository
}

func NewUserUsecase(userRepository UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
	}
}

func (userUsecase *UserUsecase) Signup(ctx context.Context, user *domain.User) (string, error) {
	user.LastVisitAt = time.Now()
	user.RegisteredAt = time.Now()

	id, err := userUsecase.userRepository.CreateUser(ctx, user)
	if err != nil {
		return utils.EmptyString, err
	}

	return id, nil
}

func (userUsecase *UserUsecase) GetUserInfo(ctx context.Context, id string) (*domain.User, error) {
	user, err := userUsecase.userRepository.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userUsecase *UserUsecase) GetAllUsersInfo(ctx context.Context) (*[]domain.User, error) {
	users, err := userUsecase.userRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (userUsecase *UserUsecase) UpdateUserInfo(ctx context.Context, user *domain.User) error {
	err := userUsecase.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (userUsecase *UserUsecase) DeleteUserAccount(ctx context.Context, id string) error {
	err := userUsecase.userRepository.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
