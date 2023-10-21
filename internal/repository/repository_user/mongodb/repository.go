package repository_user

import (
	"context"
	"time"

	"github.com/Meystergod/gochat/internal/apperror"
	"github.com/Meystergod/gochat/internal/domain"
	"github.com/Meystergod/gochat/internal/utils"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(storage *mongo.Database, collection string) *UserRepository {
	return &UserRepository{
		collection: storage.Collection(collection),
	}
}

func (userRepository *UserRepository) CreateUser(ctx context.Context, domainUser *domain.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	repositoryUser, err := userToRepository(domainUser, MethodCreate)
	if err != nil {
		err = errors.Wrap(err, "failed to convert user model")
		return utils.EmptyString, apperror.NewAppError(apperror.ErrorConvertModel, err.Error())
	}

	result, err := userRepository.collection.InsertOne(ctx, repositoryUser)
	if err != nil {
		err = errors.Wrap(err, "failed to create user")
		return utils.EmptyString, apperror.NewAppError(apperror.ErrorCreateOne, err.Error())
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		err = errors.Wrap(err, "failed to convert user id to oid")
		return utils.EmptyString, apperror.NewAppError(apperror.ErrorConvert, err.Error())
	}

	return oid.Hex(), nil
}

func (userRepository *UserRepository) GetUser(ctx context.Context, id string) (*domain.User, error) {
	var repositoryUser *User

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		err = errors.Wrap(err, "failed to convert user id to oid")
		return nil, apperror.NewAppError(apperror.ErrorConvert, err.Error())
	}

	filter := bson.M{"_id": oid}

	result := userRepository.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		err = errors.Wrap(result.Err(), "failed to get user")
		return nil, apperror.NewAppError(apperror.ErrorGetOne, err.Error())
	}

	if err = result.Decode(&repositoryUser); err != nil {
		err = errors.Wrap(err, "failed to decode user mongo object to struct")
		return nil, apperror.NewAppError(apperror.ErrorDecode, err.Error())
	}

	domainUser := userToDomain(repositoryUser)

	return &domainUser, nil
}

func (userRepository *UserRepository) GetAllUsers(ctx context.Context) (*[]domain.User, error) {
	var repositoryUsers []User

	filter := bson.M{}

	cursor, err := userRepository.collection.Find(ctx, filter)
	if err != nil {
		err = errors.Wrap(err, "failed to get all users")
		return nil, apperror.NewAppError(apperror.ErrorGetAll, err.Error())
	}

	if err = cursor.All(ctx, &repositoryUsers); err != nil {
		err = errors.Wrap(err, "failed to decode all users mongo objects to struct")
		return nil, apperror.NewAppError(apperror.ErrorDecode, err.Error())
	}

	var domainUsers []domain.User

	for _, u := range repositoryUsers {
		domainUsers = append(domainUsers, userToDomain(&u))
	}

	return &domainUsers, nil
}

func (userRepository *UserRepository) UpdateUser(ctx context.Context, domainUser *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	repositoryUser, err := userToRepository(domainUser, MethodUpdate)
	if err != nil {
		err = errors.Wrap(err, "failed to convert user model")
		return apperror.NewAppError(apperror.ErrorConvertModel, err.Error())
	}

	userByte, err := bson.Marshal(&repositoryUser)
	if err != nil {
		err = errors.Wrap(err, "failed to marshal user model to bytes")
		return apperror.NewAppError(apperror.ErrorDecode, err.Error())
	}

	var object bson.M

	err = bson.Unmarshal(userByte, &object)
	if err != nil {
		err = errors.Wrap(err, "failed to unmarshal bytes to bson")
		return apperror.NewAppError(apperror.ErrorDecode, err.Error())
	}

	delete(object, "_id")

	update := bson.M{
		"$set": object,
	}

	filter := bson.M{"_id": repositoryUser.ID}

	result, err := userRepository.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		err = errors.Wrap(err, "failed to update user")
		return apperror.NewAppError(apperror.ErrorUpdateOne, err.Error())
	}

	if result.MatchedCount == 0 {
		err = errors.New("can not be matched: failed to get user in database for update")
		return apperror.NewAppError(apperror.ErrorUpdateOne, err.Error())
	}

	return nil
}

func (userRepository *UserRepository) DeleteUser(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		err = errors.Wrap(err, "failed to convert user id to oid")
		return apperror.NewAppError(apperror.ErrorConvert, err.Error())
	}

	filter := bson.M{"_id": oid}

	result, err := userRepository.collection.DeleteOne(ctx, filter)
	if err != nil {
		err = errors.Wrap(err, "failed to delete user")
		return apperror.NewAppError(apperror.ErrorDeleteOne, err.Error())
	}

	if result.DeletedCount == 0 {
		err = errors.New("can not be deleted: failed to get user in database for delete")
		return apperror.NewAppError(apperror.ErrorDeleteOne, err.Error())
	}

	return nil
}
